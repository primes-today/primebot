package primebot

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
)

type Poster interface {
	Post(context.Context, *big.Int) error
}

func NewWriterPoster(w io.Writer) *WritePoster {
	return &WritePoster{w: w}
}

func NewStdoutPoster() *WritePoster {
	return &WritePoster{w: os.Stdout}
}

type WritePoster struct {
	w io.Writer
}

func (w *WritePoster) Post(ctx context.Context, status *big.Int) error {
	_, err := fmt.Fprintf(w.w, "%s\n", status.Text(10))
	return err
}
