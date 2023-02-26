package primebot

import (
	"context"
	"testing"
)

// 100 years of hourly postings, ignoring leap years
const MAX_PRIMES = 100 * 365 * 24

func TestGeneratorsIdentical(t *testing.T) {
	ctx := context.Background()

	td := NewTrialDivisionGenerator(0)
	pp := NewProbablyPrimeGenerator(0)

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

		if p1 != p2 {
			t.Errorf("got unequal primes; p1: %d, p2: %d", p1, p2)
		}

		count = count + 1
		if count >= MAX_PRIMES {
			// successfully reached count of equal primes
			return
		}
	}
}
