package parser

import (
	"fmt"
	"sync"
	"time"

	"github.com/tomarrell/lbadd/internal/parser/ast"
	"github.com/tomarrell/lbadd/internal/parser/scanner"
	"github.com/tomarrell/lbadd/internal/parser/scanner/token"
	"go.uber.org/multierr"
)

type isolate struct {
	sql string

	stream token.Stream

	pos    int
	tokens []token.Token // should be removed once parser and scanner actually work concurrently

	errorsLock sync.Mutex
	errors     []error
}

func newIsolate(sql string) *isolate {
	return &isolate{
		sql: sql,

		stream: token.NewStream(),
	}
}

func (i *isolate) Parse() (*ast.Query, error) {
	scanner := scanner.New([]rune(i.sql), i.stream)
	go scanner.Scan()
	return i.parse()
}

func (i *isolate) parse() (q *ast.Query, err error) {
	// handle package-internal panics
	defer func() {
		if recovered := recover(); recovered != nil {
			if _err, ok := recovered.(error); ok {
				// if an actual error was panicked, return it together with all
				// parse errors that the isolate encountered
				err = fmt.Errorf("fatal error: %w (other errors while parsing: %v)", _err, i.errors)
			} else {
				// if anything but an error was panicked, return it as an
				// ErrUnrecoverable together with all parse errors that the
				// isolate encountered
				err = fmt.Errorf("fatal error: %w: %v (other errors while parsing: %v)", ErrUnrecoverable, recovered, i.errors)
			}
		}
	}()

	// This should be removed once parser and scanner actually work
	// concurrently. The 5 second timeout is a fail-safe that returns a timeout
	// error if the scanner does not emit all tokens within 5 seconds. This
	// should be plenty of time. If 5 seconds is not enough, this will also fail
	// with the same message. The fail-safe is meant for the case that the
	// scanner is caught in an infinite state loop, meaning that one or more
	// states were poorly defined and the scanner is stuck in a infinite
	// recursion (concerning the states). Because the states are executed in a
	// loop and return the next state, the recursion is not an actual recursion
	// and will not abort with a stackoverflow error. This is why this failsafe
	// is needed. The concurrent implementation, that is desired, will accept a
	// context, that can optionally provide a deadline to the scanner/parser.
	select {
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("accumulate tokens: %w (5sec)", ErrTimeout)
	case <-i.accumulateTokens():
	}

	q = &ast.Query{
		Node:            ast.NewNodeInfo(1, 1, 0, len(i.sql), ast.Root, i.sql),
		QueryExpression: i.parseQueryExpression(),
	}

	if len(i.errors) > 0 {
		err = fmt.Errorf("%w: %v", ErrParserError, multierr.Combine(i.errors...))
	}

	return
}

// accumulateTokens reads all tokens from the stream and stores them inside this
// isolate. This makes lookaheads and lookbacks more easy than implementing
// indefinite caching.
func (i *isolate) accumulateTokens() <-chan struct{} { // should be removed once parser and scanner actually work concurrently
	ch := make(chan struct{})
	go func() {
		for t := i.stream.Take(); t.Type() != token.EOF; t = i.stream.Take() { // EOF token is omitted
			if t.Type() == token.Error {
				i.foundError(fmt.Errorf("%w: %v", ErrScannerError, t.Value()))
			} else {
				i.tokens = append(i.tokens, t)
			}
		}
		close(ch)
	}()
	return ch
}

// foundError adds the given error to this isolate's error collection. This
// method is safe for concurrent use.
func (i *isolate) foundError(err error) {
	i.errorsLock.Lock()
	defer i.errorsLock.Unlock()
	i.errors = append(i.errors, err)
}

// next returns the next token from the token stream, or nil if there are no
// more. ok=false implies, that the token stream has reached its end, and there
// is no and there will be no more tokens to read. ok=false also implies, that
// the current position of the parser on the token stream has not shifted in any
// direction.
func (i *isolate) next() (t token.Token, ok bool) {
	if i.done() {
		return nil, false
	}

	next := i.tokens[i.pos]
	i.pos++

	return next, true
}

func (i *isolate) peek() (token.Token, bool) {
	if i.done() {
		return nil, false
	}
	return i.tokens[i.pos], true
}

func (i *isolate) goback() {
	i.pos--
}

func (i *isolate) done() bool {
	return i.pos >= len(i.tokens)
}

func (i *isolate) acceptToken(typ token.Type) (token.Token, bool) {
	if next, ok := i.next(); ok {
		return next, ok
	}
	i.goback()
	return nil, false
}

func (i *isolate) acceptTokenSequence(types ...token.Type) ([]token.Token, bool) {
	pos := i.pos
	result := make([]token.Token, len(types))

	for _, t := range types {
		next, ok := i.acceptToken(t)
		if !ok {
			i.pos = pos // reset position
			return result, false
		}
		result = append(result, next)
	}

	return result, true
}

func (i *isolate) acceptOneOfTypes(types ...token.Type) (token.Token, bool) {
	next, ok := i.next()
	if !ok {
		return nil, false
	}

	for _, t := range types {
		if t == next.Type() {
			return next, true
		}
	}
	i.goback()
	return nil, false
}
