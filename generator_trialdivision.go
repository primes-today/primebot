package primebot

import (
	"context"
	"math/big"
	"sync"
)

func NewTrialDivisionGenerator(start *big.Int) *TrialDivisionGenerator {
	cur := &big.Int{}
	cur.Set(start)
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)

	return &TrialDivisionGenerator{
		cur:   cur,
		zero:  zero,
		one:   one,
		two:   two,
		mutex: &sync.Mutex{},
	}
}

type TrialDivisionGenerator struct {
	cur   *big.Int
	zero  *big.Int
	one   *big.Int
	two   *big.Int
	mutex *sync.Mutex
}

func (t *TrialDivisionGenerator) SetStart(start *big.Int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	cur := &big.Int{}
	cur.Set(start)
	t.cur = cur
}

func (t *TrialDivisionGenerator) Generate(ctx context.Context) (*big.Int, error) {
	defer t.cur.Add(t.cur, t.one)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	ret := &big.Int{}
	trial := big.NewInt(2)
	mod := &big.Int{}
	max := &big.Int{}

	// handle == 2
	if t.cur.Cmp(trial) == 0 {
		return ret.Set(t.cur), nil
	}
	// handle < 2
	if t.cur.Cmp(trial) < 0 {
		t.cur.Set(trial)
	}

	for {
		select {
		case <-ctx.Done():
			return ret, ErrCanceled
		default:
			trial.SetUint64(2)
			first := true
			max.Sqrt(t.cur)

		OUTER:
			for {
				select {
				case <-ctx.Done():
					return ret, ErrCanceled
				default:
					if trial.Cmp(max) > 0 {
						return ret.Set(t.cur), nil
					}
					if mod.Mod(t.cur, trial).Cmp(t.zero) == 0 {
						break OUTER
					}
					if first {
						trial.Add(trial, t.one)
						first = false
					} else {
						trial.Add(trial, t.two)
					}
				}
			}

			t.cur.Add(t.cur, t.one)
		}
	}
}
