package primebot

import (
	"context"
	"testing"
)

// 50 years of hourly postings, ignoring leap years
const MAX_PRIMES = 50 * 365 * 24

func TestGeneratorsIdentical(t *testing.T) {
	ctx := context.Background()

	td := NewTrialDivisionGenerator(0)
	pp := NewProbablyPrimeGenerator(0)

	count := 0
	for {
		p1, err1 := td.Generate(ctx)
		p2, err2 := pp.Generate(ctx)
		if err1 != nil {
			if err1 == err2 && err1 == ErrOverflow {
				// success
				return
			}

			t.Fatal(err1)
		}
		if err2 != nil {
			t.Fatal(err2)
		}

		if p1 != p2 {
			t.Errorf("got unequal primes; p1: %d, p2: %d", p1, p2)
		}

		count = count + 1
		if count >= MAX_PRIMES {
			// done
			return
		}
	}
}
