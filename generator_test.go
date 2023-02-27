package primebot

import (
	"context"
	"math/big"
	"testing"
)

// 100 years of hourly postings, ignoring leap years
const MAX_PRIMES = 100 * 365 * 24

func TestGeneratorsIdentical(t *testing.T) {
	ctx := context.Background()

	td := NewTrialDivisionGenerator(big.NewInt(0))
	pp := NewProbablyPrimeGenerator(big.NewInt(0))

	count := 0
	for {
		p1, err := td.Generate(ctx)
		if err != nil {
			t.Fatal(err)
		}
		p2, err := pp.Generate(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if p1.Cmp(p2) != 0 {
			t.Errorf("got unequal primes; p1: %d, p2: %d", p1, p2)
		}

		count = count + 1
		if count >= MAX_PRIMES {
			// successfully reached count of equal primes
			return
		}
	}
}
