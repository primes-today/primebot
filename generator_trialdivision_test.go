package primebot

import (
	"context"
	"math/big"
	"testing"
)

func TestTrialDivisionGenerator(t *testing.T) {
	ctx := context.Background()
	cases := map[int64][]uint64{
		0: {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31},
		5: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
		6: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
	}

	for start, expected := range cases {
		g := NewTrialDivisionGenerator(big.NewInt(start))
		for _, ex := range expected {
			p, err := g.Generate(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if p.Uint64() != ex {
				t.Errorf("expected equal primes, got %v, wanted %v", p, ex)
			}
		}
	}
}

func TestTrialDivisionGeneratorCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	start := &big.Int{}
	start.SetUint64(uint64(1 << 32))

	g := NewTrialDivisionGenerator(start)
	go func() {
		_, err := g.Generate(ctx)
		if err != ErrCanceled {
			t.Errorf("expected cancellation error, got %v", err)
		}
	}()
	cancel()
}

func TestTrialDivisionGeneratorDoesNotOverflow(t *testing.T) {
	ctx := context.Background()
	start := &big.Int{}
	start.SetUint64(uint64((1 << 64) - 1))

	expected := &big.Int{}
	expected.SetString("18446744073709551629", 10)

	g := NewTrialDivisionGenerator(start)
	n, err := g.Generate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if n.Cmp(expected) != 0 {
		t.Errorf("expected prime %s but got %s", expected, n)
	}
}
