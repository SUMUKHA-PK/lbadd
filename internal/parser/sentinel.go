package parser

var _ error = (*sentinel)(nil) // ensure that sentinel implements error

// sentinel is a helper type that implements the error interface to allow for
// constant errors.
type sentinel string

func (s sentinel) Error() string { return string(s) }
