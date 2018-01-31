package primebot

import (
	"context"
	"math/big"
	"sync"
)

// NewProbablyPrimeGenerator creates a new ProbablyPrimeGenerator, starting at
// the given number
func NewProbablyPrimeGenerator(start uint64) *ProbablyPrimeGenerator {
	cur := &big.Int{}
	cur.SetUint64(start)
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

func (p *ProbablyPrimeGenerator) SetStart(start uint64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	cur := &big.Int{}
	cur.SetUint64(start)
	p.cur = cur
}

// Generate generates the next prime number from the ProbablyPrime generator
func (p *ProbablyPrimeGenerator) Generate(ctx context.Context) (uint64, error) {
	defer p.cur.Add(p.cur, p.one)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for {
		select {
		case <-ctx.Done():
			return 0, ErrCanceled
		default:
			if !p.cur.IsUint64() {
				return 0, ErrOverflow
			}
			if p.cur.ProbablyPrime(1) {
				return p.cur.Uint64(), nil
			}
			p.cur.Add(p.cur, p.one)
		}
	}
}
