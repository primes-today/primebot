package primebot

import (
	"context"
	"testing"
)

func TestProbablyPrimeGenerator(t *testing.T) {
	ctx := context.Background()
	cases := map[uint64][]uint64{
		0: []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31},
		7: []uint64{7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
		6: []uint64{7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
	}

	for start, expected := range cases {
		g := NewProbablyPrimeGenerator(start)
		for _, ex := range expected {
			p, err := g.Generate(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if p != ex {
				t.Errorf("expected equal primes, got %v, wanted %v", p, ex)
			}
		}
	}
}

func TestProbablyPrimeGeneratorCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	start := uint64(1 << 32)

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
	start := uint64((1 << 64) - 1)

	g := NewProbablyPrimeGenerator(start)
	n, err := g.Generate(ctx)
	if err != ErrOverflow {
		t.Errorf("expected overflow error, got %v, %v", n, err)
	}
}
