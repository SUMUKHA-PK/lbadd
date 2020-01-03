package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomarrell/lbadd/internal/parser/scanner/token"
)

type testStream struct {
	tokens []token.Token
	closed bool
}

func (s *testStream) Peek() token.Token {
	if s.closed {
		panic("already closed")
	}
	return s.tokens[0]
}

func (s *testStream) Take() token.Token {
	if s.closed {
		panic("already closed")
	}
	t := s.tokens[0]
	s.tokens[0] = nil
	s.tokens = s.tokens[1:]
	return t
}

func (s *testStream) Push(t token.Token) {
	if s.closed {
		panic("already closed")
	}
	s.tokens = append(s.tokens, t)
}

func (s *testStream) Close() error {
	s.closed = true
	return nil
}

func (s *testStream) Closed() bool {
	return s.closed
}

func (s *testStream) errors() (res []token.Token) {
	for _, t := range s.tokens {
		if t.Type() == token.Error {
			res = append(res, t)
		}
	}
	return
}

type testScanner struct {
	*Scanner
	t *testing.T
	a *assert.Assertions
}

func newTestScanner(t *testing.T, input string, st token.Stream) *testScanner {
	return &testScanner{
		Scanner: New([]rune(input), st),
		t:       t,
		a:       assert.New(t),
	}
}

func (s *testScanner) executeState(st state) {
	s.current = st
	s.executeCurrentState()
}

func (s *testScanner) assertCurrentStateIsNil() {
	s.a.Nil(s.current, "current state must be nil, but was not, meaning that an error occurred (current state: '%s')", s.current)
}
