package primebot

import (
	"context"
	"testing"
)

func TestTrialDivisionGenerator(t *testing.T) {
	ctx := context.Background()
	cases := map[uint64][]uint64{
		0: {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31},
		7: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
		6: {7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43},
	}

	for start, expected := range cases {
		g := NewTrialDivisionGenerator(start)
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

func TestTrialDivisionGeneratorCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	start := uint64(1 << 32)

	g := NewTrialDivisionGenerator(start)
	go func() {
		_, err := g.Generate(ctx)
		if err != ErrCanceled {
			t.Errorf("expected cancellation error, got %v", err)
		}
	}()
	cancel()
}

func TestTrialDivisionGeneratorOverflow(t *testing.T) {
	ctx := context.Background()
	start := uint64((1 << 64) - 1)

	g := NewTrialDivisionGenerator(start)
	n, err := g.Generate(ctx)
	if err != ErrOverflow {
		t.Errorf("expected overflow error, got %v, %v", n, err)
	}
}
