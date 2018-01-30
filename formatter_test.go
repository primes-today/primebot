package primebot

import (
	"testing"
)

func TestPlainFormat(t *testing.T) {
	fmt := PlainFormat{}
	if s := fmt.Format(10001); s != "10001" {
		t.Errorf("got unexpected string %v", s)
	}
}
