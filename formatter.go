package primebot

import (
	"fmt"
)

// Formatter defines an interface for formatting integers to displayable strings
type Formatter interface {
	Format(uint64) string
}

// PlainFormat is a formatter which outputs integers as a string value, with no
// additional formatting
type PlainFormat struct{}

// Format formats a uint to a string
func (p PlainFormat) Format(n uint64) string {
	return fmt.Sprintf("%v", n)
}
