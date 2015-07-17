package gobackoff

import (
	"math/rand"
	"time"

	"golang.org/x/net/context"
)

const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.8
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 300 * time.Second
	DefaultMaxElapsedTime      = 900 * time.Second
)

var (
	DefaultParam = BackOffParams{
		InitialInterval:     DefaultInitialInterval,
		RandomizationFactor: DefaultRandomizationFactor,
		Multiplier:          DefaultMultiplier,
		MaxInterval:         DefaultMaxInterval,
		MaxElapsedTime:      DefaultMaxElapsedTime,
	}
)

type BackOffParams struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

type BackOff struct {
	Ctx context.Context
	BackOffParams
	startTime       time.Time
	currentInterval time.Duration
}

func NewBackoff() *BackOff {
	ctx, _ := context.WithCancel(context.Background())
	return &BackOff{
		Ctx:           ctx,
		BackOffParams: DefaultParam,
	}
}
func NewBackoffParam(ctx context.Context, p BackOffParams) *BackOff {
	//ctx, _ := context.WithCancel(context.Background())
	return &BackOff{
		Ctx:           ctx,
		BackOffParams: p,
	}
}

func (b *BackOff) Retry(cb func() error) error {
	var err error
	var next time.Duration
	var stop bool
	b.startTime = time.Now()
	b.currentInterval = b.InitialInterval

	for {
		if err = cb(); err == nil {
			return nil
		}

		if next, stop = b.Next(b.currentInterval); stop {
			return err
		}

		select {
		case <-b.Ctx.Done():
			return b.Ctx.Err()
		case <-time.After(next):
		}
	}

}

func (b *BackOff) Next(current time.Duration) (next time.Duration, stop bool) {
	if b.MaxElapsedTime != -1 && time.Now().Sub(b.startTime) > b.MaxElapsedTime {
		return next, true
	}
	delta := float64(current) * b.RandomizationFactor
	minInterval := float64(current) - delta
	maxInterval := float64(current) + delta
	next = time.Duration(minInterval + (rand.Float64() * (maxInterval - minInterval + 1)))

	if float64(next) >= float64(b.MaxInterval)/b.Multiplier {
		next = b.MaxInterval
	} else {
		next = time.Duration(float64(next) * b.Multiplier)
	}
	return next, false
}
