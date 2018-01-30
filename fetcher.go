package primebot

import (
	"math/rand"
	"time"
)

type Status struct {
	num uint64
	ts  time.Time
}

type Fetcher interface {
	Fetch() (*Status, error)
}

func NewRandFetcher() *RandFetcher {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &RandFetcher{r: r}
}

type RandFetcher struct {
	r *rand.Rand
}

func (r *RandFetcher) Fetch() (*Status, error) {
	now := time.Now()
	return &Status{
		num: r.r.Uint64(),
		ts:  now.Add(-10 * time.Second),
	}, nil
}
