package scanner

type SyntaxError struct {
	offset    int
	line, col int

	message string
}
