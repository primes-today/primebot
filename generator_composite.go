package primebot

import (
	"context"
	"errors"
	"math/big"
	"sync"
)

func NewCompositeGenerator(start *big.Int) *CompositeGenerator {
	maxP := &big.Int{}
	maxP.SetUint64(uint64((1 << 64) - 1))
	gen := []Generator{
		NewProbablyPrimeGenerator(start),
		NewTrialDivisionGenerator(start),
	}

	return &CompositeGenerator{
		mutex: &sync.Mutex{},
		one:   big.NewInt(1),
		gen:   gen,
	}
}

type CompositeGenerator struct {
	mutex *sync.Mutex
	one   *big.Int
	gen   []Generator
}

func (c *CompositeGenerator) SetStart(start *big.Int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, gen := range c.gen {
		gen.SetStart(start)
	}
}

func (c *CompositeGenerator) Generate(ctx context.Context) (*big.Int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, gen := range c.gen {
		n, err := gen.Generate(ctx)
		if err == ErrOverflow {
			continue
		}

		for _, gen := range c.gen {
			gen.SetStart(n)
		}

		return n, err
	}

	return &big.Int{}, errors.New("failed to generate prime; all generators fell through")
}
