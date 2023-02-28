package primebot

import (
	"context"
	"math/big"
	"testing"
)

func TestProbablyPrimeGenerator(t *testing.T) {
	ctx := context.Background()
	cases := map[int64][]uint64{
		0: {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31},
		5: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
		6: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
	}

	for start, expected := range cases {
		g := NewProbablyPrimeGenerator(big.NewInt(start))
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

func TestProbablyPrimeGeneratorCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	start := &big.Int{}
	start.SetUint64(uint64(1 << 32))

	g := NewProbablyPrimeGenerator(start)
	go func() {
		_, err := g.Generate(ctx)
		if err != ErrCanceled {
			t.Errorf("expected cancellation error, got %v", err)
		}
	}()
	cancel()
}

func TestProbablyPrimeGeneratorOverflow(t *testing.T) {
	ctx := context.Background()
	start := &big.Int{}
	start.SetUint64(uint64((1 << 64) - 1))

	g := NewProbablyPrimeGenerator(start)
	n, err := g.Generate(ctx)
	if err != ErrOverflow {
		t.Errorf("expected overflow error, got %v, %v", n, err)
	}
}
