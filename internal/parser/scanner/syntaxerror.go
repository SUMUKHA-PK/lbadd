package scanner

import "fmt"

type SyntaxError struct {
	offset    int
	line, col int

	message string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at %d:%d (offset %d): %s", e.line, e.col, e.offset, e.message)
}
