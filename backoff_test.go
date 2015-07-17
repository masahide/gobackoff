package gobackoff

import (
	"fmt"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	param := BackOffParams{
		InitialInterval:     100 * time.Millisecond,
		RandomizationFactor: 0.8,
		Multiplier:          1.5,
		MaxInterval:         500 * time.Millisecond,
		MaxElapsedTime:      1000 * time.Millisecond,
	}
	b := NewBackoff(param)

	start := time.Now()
	times := []time.Duration{}
	err := b.Retry(func() error {
		times = append(times, time.Now().Sub(start))
		return fmt.Errorf("ng")
	})

	for i, t := range times {
		fmt.Printf("%02d,%s\n", i, t)
	}
	if err == nil {
		t.Errorf("got err:nil\ntimes %#v", err, times)
	}
	max := param.MaxElapsedTime + time.Duration(float64(param.InitialInterval)*param.Multiplier)
	last := times[len(times)-1]
	if last > max {
		t.Errorf("got last:%s   want : <%s", last, max)
	}
}

func TestRetryCancel(t *testing.T) {
	param := BackOffParams{
		InitialInterval:     100 * time.Millisecond,
		RandomizationFactor: 0.8,
		Multiplier:          1.5,
		MaxInterval:         500 * time.Millisecond,
		MaxElapsedTime:      1000 * time.Millisecond,
	}
	b := NewBackoff(param)

	start := time.Now()
	times := []time.Duration{}
	err := b.Retry(func() error {
		times = append(times, time.Now().Sub(start))
		return fmt.Errorf("ng")
	})

	for i, t := range times {
		fmt.Printf("%02d,%s\n", i, t)
	}
	if err == nil {
		t.Errorf("got err:nil\ntimes %#v", err, times)
	}
	max := param.MaxElapsedTime + time.Duration(float64(param.InitialInterval)*param.Multiplier)
	last := times[len(times)-1]
	if last > max {
		t.Errorf("got last:%s   want : <%s", last, max)
	}
}
