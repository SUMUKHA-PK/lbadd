package scanner

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/tomarrell/lbadd/internal/parser/scanner/matcher"
	"github.com/tomarrell/lbadd/internal/parser/scanner/token"
)

// state is a type alias for a function that takes a scanner and returns another
// state. Such functions (or states) will be invoked by the scanner. It will
// pass itself as the argument, and the returned state will be chained and
// executed next. If a state evaluates to nil, the initial state will be
// chained, if the scanner has not reached the EOF yet.
//
// It is a state's responsibility to reset the scanner's position etc. to the
// state they were in, when the state was entered, if the state does not
// evaluate to nil. A state can do this by using checkpoints as shown in the
// following example.
//
//	func myState(s *Scanner) state {
//		chck := s.checkpoint()
//		if s.accept(MyMatcher) {
//			// maybe emit a token here
//			return nil
//		}
//		s.restore(chck) // resets start, pos, line, col and the other position attributes
//		return errorf("myerror")
//	}
type state func(*Scanner) state

func (st state) String() string {
	return runtime.FuncForPC(reflect.ValueOf(st).Pointer()).Name()
}

type checkpoint struct {
	start int
	pos   int

	startLine, startCol int
	line, lastCol, col  int
}

type Scanner struct {
	input []rune
	start int
	pos   int

	current state
	stream  token.Stream

	startLine, startCol int
	line, lastCol, col  int

	closed bool
	doneCh chan struct{}
}

func New(input []rune, stream token.Stream) *Scanner {
	return &Scanner{
		input: input,
		start: 0,
		pos:   0,

		current: initial,
		stream:  stream,

		startLine: 1,
		startCol:  1,
		line:      1, // line starts at 1, because it should be human readable and editor line and column numbers usually start at 1
		lastCol:   1,
		col:       1, // col starts at 1, because it should be human readable and editor line and column numbers usually start at 1

		doneCh: make(chan struct{}),
	}
}

func (s *Scanner) Scan() {
	// EOF must be emitted after recovering from a crash, so that the parser
	// doesn't miss the error
	defer s.emit(token.EOF)

	// recover from syntax errors that cannot be handled before EOF is emitted
	defer func() {
		if recovered := recover(); recovered != nil {
			if err, ok := recovered.(SyntaxError); ok {
				s.stream.Push(token.New(s.line, s.col, s.pos, s.pos-s.start, token.Error, fmt.Sprintf("recovered: %v", err)))
			} else {
				panic(recovered) // re-panic if it's not a syntax error
			}
		}
	}()

	for !s.done() {
		s.executeCurrentState()
		if s.current == nil {
			s.current = initial
		}
	}

	close(s.doneCh)
}

func (s *Scanner) Done() <-chan struct{} {
	return s.doneCh
}

// Close will cause this scanner to not execute any more states and emit an EOF
// token.
func (s *Scanner) Close() error {
	s.closed = true
	return nil
}

func (s *Scanner) executeCurrentState() {
	s.current = s.current(s)
}

func (s *Scanner) done() bool {
	return s.closed ||
		s.pos >= len(s.input)
}

func (s *Scanner) next() rune {
	if s.done() {
		panic(SyntaxError{
			offset:  s.pos,
			line:    s.line,
			col:     s.col,
			message: fmt.Sprintf("state '%v' tried to read another rune, but scanner already reached EOF at offset %d (%d:%d)", s.current, s.pos, s.line, s.col),
		})
	}

	next := s.input[s.pos]

	// update line and column information
	if next == '\n' { // TODO if next is line delimiter
		s.line++
		s.lastCol = s.col
		s.col = 1
	} else {
		s.col++
	}

	// update current scanner position
	s.pos++

	return next
}

func (s *Scanner) peek() rune {
	return s.input[s.pos]
}

func (s *Scanner) goback() {
	s.pos--

	// update line and column information
	if s.col == 1 {
		s.line--
		s.col = s.lastCol
	} else {
		s.col--
	}
}

func (s *Scanner) checkpoint() checkpoint {
	return checkpoint{
		start:     s.start,
		pos:       s.pos,
		startLine: s.startLine,
		startCol:  s.startCol,
		line:      s.line,
		lastCol:   s.lastCol,
		col:       s.col,
	}
}

func (s *Scanner) restore(chck checkpoint) {
	s.start = chck.start
	s.pos = chck.pos
	s.startLine = chck.startLine
	s.startCol = chck.startCol
	s.line = chck.line
	s.lastCol = chck.lastCol
	s.col = chck.col
}

func (s *Scanner) emit(t token.Type) {
	tok := token.New(s.line, s.col, s.start, s.pos-s.start, t, string(s.input[s.start:s.pos]))
	s.stream.Push(tok)

	s.start = s.pos
	s.startLine = s.line
	s.startCol = s.col
}

func (s *Scanner) accept(m matcher.M) bool {
	if m.Matches(s.next()) {
		return true
	}
	s.goback()
	return false
}

func (s *Scanner) acceptMultiple(m matcher.M) (matched uint) {
	for s.accept(m) {
		matched++
	}
	return
}

func (s *Scanner) acceptString(str string) bool {
	chck := s.checkpoint()
	for _, r := range str {
		if r != s.next() {
			s.restore(chck)
			return false
		}
	}
	return true
}

func (s *Scanner) peekString(str string) bool {
	chck := s.checkpoint()
	defer s.restore(chck)

	for _, r := range str {
		if r != s.next() {
			return false
		}
	}
	return true
}
