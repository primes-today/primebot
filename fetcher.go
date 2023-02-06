package primebot

import (
	"context"
	"math/rand"
	"time"
)

type Status struct {
	LastStatus uint64
	Posted     time.Time
}

type Fetcher interface {
	Fetch(context.Context) (*Status, error)
}

func NewRandFetcher() *RandFetcher {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return &RandFetcher{r: r}
}

type RandFetcher struct {
	r *rand.Rand
}

func (r *RandFetcher) Fetch(ctx context.Context) (*Status, error) {
	now := time.Now()
	return &Status{
		LastStatus: r.r.Uint64(),
		Posted:     now.Add(-10 * time.Second),
	}, nil
}
