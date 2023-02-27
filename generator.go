package primebot

import (
	"context"
	"errors"
	"math/big"
)

var (
	// ErrCanceled is returned from a generator when its context is cancelled
	ErrCanceled = errors.New("cancelled")
	// ErrOverflow is returned from a generator when the generator number type's
	// max is overflowed
	ErrOverflow = errors.New("overflow")
)

type Generator interface {
	SetStart(*big.Int)
	Generate(context.Context) (*big.Int, error)
}
