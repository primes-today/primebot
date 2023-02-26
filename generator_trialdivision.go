package primebot

import (
	"context"
	"math/big"
	"sync"
)

func NewTrialDivisionGenerator(start uint64) *TrialDivisionGenerator {
	cur := &big.Int{}
	cur.SetUint64(start)
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

func (t *TrialDivisionGenerator) SetStart(start uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	cur := &big.Int{}
	cur.SetUint64(start)
	t.cur = cur
}

func (t *TrialDivisionGenerator) Generate(ctx context.Context) (uint64, error) {
	trial := big.NewInt(2)
	mod := &big.Int{}
	max := &big.Int{}

	defer t.cur.Add(t.cur, t.one)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// handle == 2
	if t.cur.Cmp(trial) == 0 {
		return t.cur.Uint64(), nil
	}
	// handle < 2
	if t.cur.Cmp(trial) < 0 {
		t.cur.Set(trial)
	}

	for {
		select {
		case <-ctx.Done():
			return 0, ErrCanceled
		default:
			if !t.cur.IsUint64() {
				return 0, ErrOverflow
			}

			trial.SetUint64(2)
			first := true
			max.Sqrt(t.cur)

		OUTER:
			for {
				select {
				case <-ctx.Done():
					return 0, ErrCanceled
				default:
					if trial.Cmp(max) > 0 {
						return t.cur.Uint64(), nil
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
