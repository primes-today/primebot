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
	three := big.NewInt(3)
	six := big.NewInt(6)

	return &TrialDivisionGenerator{
		cur:   cur,
		zero:  zero,
		one:   one,
		two:   two,
		three: three,
		six:   six,
		mutex: &sync.Mutex{},
	}
}

type TrialDivisionGenerator struct {
	cur   *big.Int
	zero  *big.Int
	one   *big.Int
	two   *big.Int
	three *big.Int
	six   *big.Int
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
	t.mutex.Lock()
	defer t.mutex.Unlock()

	ret := &big.Int{}
	mod := &big.Int{}
	max := &big.Int{}
	t1 := &big.Int{}
	t2 := &big.Int{}

	for {
		t.cur.Add(t.cur, t.one)

		select {
		case <-ctx.Done():
			return ret, ErrCanceled
		default:
		}

		if t.cur.Cmp(t.two) == 0 || t.cur.Cmp(t.three) == 0 {
			return ret.Set(t.cur), nil
		}
		if t.cur.Cmp(t.one) < 0 || t.cur.Cmp(t.one) == 0 {
			continue
		}
		if mod.Mod(t.cur, t.two).Cmp(t.zero) == 0 {
			continue
		}
		if mod.Mod(t.cur, t.three).Cmp(t.zero) == 0 {
			continue
		}

		t1.SetUint64(5)
		t2.SetUint64(7)
		max.Sqrt(t.cur)

	OUTER:
		for {
			select {
			case <-ctx.Done():
				return ret, ErrCanceled
			default:
				if t1.Cmp(max) > 0 {
					return ret.Set(t.cur), nil
				}
				if mod.Mod(t.cur, t1).Cmp(t.zero) == 0 {
					break OUTER
				}
				if mod.Mod(t.cur, t2).Cmp(t.zero) == 0 {
					break OUTER
				}
				t1.Add(t1, t.six)
				t2.Add(t2, t.six)
			}
		}
	}
}
