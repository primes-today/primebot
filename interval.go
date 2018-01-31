package primebot

import (
	"context"
	"time"
)

type Ticker interface {
	Start(context.Context, time.Duration) chan time.Time
	Interval() time.Duration
}

func NewIntervalTicker(d time.Duration) *IntervalTicker {
	return &IntervalTicker{interval: d}
}

type IntervalTicker struct {
	interval time.Duration
}

func (i *IntervalTicker) Start(ctx context.Context, first time.Duration) (c chan time.Time) {
	c = make(chan time.Time)
	go func() {
		s := time.NewTimer(first)
		select {
		case tt := <-s.C:
			c <- tt
		case <-ctx.Done():
			s.Stop()
			return
		}

		t := time.NewTicker(i.interval)
		for {
			select {
			case tt := <-t.C:
				c <- tt
			case <-ctx.Done():
				t.Stop()
				break
			}
		}
	}()

	return c
}

func (i IntervalTicker) Interval() time.Duration {
	return i.interval
}
