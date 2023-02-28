package primebot

import (
	"context"
	"math/big"
	"math/rand"
	"time"
)

type Status struct {
	LastStatus *big.Int
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
	last := &big.Int{}

	return &Status{
		LastStatus: last.SetUint64(r.r.Uint64()),
		Posted:     now.Add(-10 * time.Second),
	}, nil
}
