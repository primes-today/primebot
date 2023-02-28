package primebot

import (
	"context"
	"math/big"
	"sync"
)

// NewProbablyPrimeGenerator creates a new ProbablyPrimeGenerator, starting at
// the given number
func NewProbablyPrimeGenerator(start *big.Int) *ProbablyPrimeGenerator {
	cur := &big.Int{}
	cur.Set(start)
	one := big.NewInt(1)

	return &ProbablyPrimeGenerator{
		cur:   cur,
		one:   one,
		mutex: &sync.Mutex{},
	}
}

// ProbablyPrimeGenerator is a prime number generator that uses the math/big
// package's ProbablyPrime; it is 100% accurate up to the max value of a uint,
// at which point it overflows; this cannot be used to generate any prime number
// larger than 2**64
type ProbablyPrimeGenerator struct {
	cur   *big.Int
	one   *big.Int
	mutex *sync.Mutex
}

func (p *ProbablyPrimeGenerator) SetStart(start *big.Int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	cur := &big.Int{}
	cur.Set(start)
	p.cur = cur
}

// Generate generates the next prime number from the ProbablyPrime generator
func (p *ProbablyPrimeGenerator) Generate(ctx context.Context) (*big.Int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.cur.Add(p.cur, p.one)

	ret := &big.Int{}
	for {
		select {
		case <-ctx.Done():
			return ret, ErrCanceled
		default:
			if !p.cur.IsUint64() {
				return ret, ErrOverflow
			}
			if p.cur.ProbablyPrime(1) {
				return ret.Set(p.cur), nil
			}
			p.cur.Add(p.cur, p.one)
		}
	}
}
