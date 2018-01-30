package primebot

import (
	"fmt"
	"io"
	"os"
)

type Poster interface {
	Post(string) error
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

func (w *WritePoster) Post(status string) error {
	_, err := fmt.Fprintf(w.w, "%v\n", status)
	return err
}
