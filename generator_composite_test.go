package primebot

import (
	"context"
	"math/big"
	"testing"
)

func TestCompositeGenerator(t *testing.T) {
	ctx := context.Background()

	start := &big.Int{}
	start.SetString("18446744073709551556", 10)

	e1 := &big.Int{}
	// last prime before 2**64 - 1
	e1.SetString("18446744073709551557", 10)
	e2 := &big.Int{}
	// first prime after 2**64 - 1
	e2.SetString("18446744073709551629", 10)

	gen := NewCompositeGenerator(start)

	p1, err := gen.Generate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if p1.Cmp(e1) != 0 {
		t.Errorf("expected %s, but got %s", e1, p1)
	}

	p2, err := gen.Generate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if p2.Cmp(e2) != 0 {
		t.Errorf("expected %s, but got %s", e2, p2)
	}
}
