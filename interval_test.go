package primebot

import (
	"context"
	"testing"
	"time"
)

func TestIntervalTicker(t *testing.T) {
	count := 4
	ctx, cancel := context.WithCancel(context.Background())
	it := NewIntervalTicker(2 * time.Second)

	now := time.Now()
	c := it.Start(ctx, 1*time.Second)

	tc := <-c
	if s := tc.Sub(now).Seconds(); s > 1.2 || s < 0.8 {
		t.Errorf("expected first tick at ~1 seconds, saw %v", s)
	}

	for {
		if count < 1 {
			cancel()
			break
		}
		now = time.Now()
		tc = <-c
		if s := tc.Sub(now).Seconds(); s > 2.2 || s < 1.8 {
			t.Errorf("expected tick at ~2 seconds, saw %v", s)
		}
		count--
	}
	close(c)
}
