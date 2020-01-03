package scanner

import (
	"fmt"

	"github.com/tomarrell/lbadd/internal/parser/scanner/token"
)

func errorf(format string, args ...interface{}) state {
	return func(s *Scanner) state {
		err := fmt.Errorf("unexpected lexical unit: %v", fmt.Sprintf(format, args...))
		errToken := token.New(s.startLine, s.startCol, s.start, s.pos-s.start, token.Error, err.Error())
		s.stream.Push(errToken)
		return nil
	}
}

func errorExpected(i interface{}) state {
	return errorf("expected %v", i)
}

func initial(s *Scanner) state {
	return scanToken
}

func scanToken(s *Scanner) state {
	// in order to avoid 3 huge switch-case statements, we merged <token>,
	// <nondelimiter token>, <delimiter token> into this very function

	next := s.peek()
	switch {
	case IdentifierStart.Matches(next):
		return scanRegularIdentifier
	case Digit.Matches(next):
		// scanUnsignedNumericLiteral OR scanLargeObjectLengthToken
	case Period.Matches(next):
		return scanUnsignedNumericLiteral
	case s.peekString(`N"`):
		return scanNationalCharacterStringLiteral
	case s.peekString(`X"`):
		return scanBinaryStringLiteral
	default:
		// handle all cases that cannot be expressed by a switch statement

		panic(SyntaxError{
			offset:  s.pos,
			line:    s.line,
			col:     s.col,
			message: fmt.Sprintf("unexpected rune '%v' near offset %v (%v:%v)", next, s.pos, s.line, s.col),
		})
	}
}

func scanLeftBracketOrTrigraph(s *Scanner) state {
	if s.accept(LeftParen) {
		return nil
	}

	return scanLeftBracketTrigraph
}

func scanLeftBracketTrigraph(s *Scanner) state {
	chck := s.checkpoint()
	if s.accept(QuestionMark) &&
		s.accept(QuestionMark) {
		if s.accept(LeftParen) {
			return nil
		}
		s.restore(chck)
		return errorExpected(LeftParen)
	}
	s.restore(chck)
	return errorExpected(QuestionMark)
}

func scanRightBracketOrTrigraph(s *Scanner) state {
	if s.accept(RightParen) {
		return nil
	}

	return scanRightBracketTrigraph
}

func scanRightBracketTrigraph(s *Scanner) state {
	chck := s.checkpoint()
	if s.accept(QuestionMark) &&
		s.accept(QuestionMark) {
		if s.accept(RightParen) {
			return nil
		}
		s.restore(chck)
		return errorExpected(RightParen)
	}
	s.restore(chck)
	return errorExpected(QuestionMark)
}

func scanRegularIdentifier(s *Scanner) state {
	return scanIdentifierBody
}

func scanIdentifierBody(s *Scanner) state {
	chck := s.checkpoint()

	if !s.accept(IdentifierStart) {
		s.restore(chck)
		return errorExpected(IdentifierStart)
	}

	_ = s.acceptMultiple(IdentifierPart)
	return nil
}

func scanLargeObjectLengthToken(s *Scanner) state {
	chck := s.checkpoint()

	// at least one Digit
	if s.acceptMultiple(Digit) == 0 {
		s.restore(chck)
		return errorExpected(Digit)
	}
	// exactly one multiplier
	if !s.accept(Multiplier) {
		s.restore(chck)
		return errorExpected(Multiplier)
	}
	return nil
}

func scanDelimitedIdentifier(s *Scanner) state {
	chck := s.checkpoint()

	if !s.accept(Quote) {
		s.restore(chck)
		return errorExpected(Quote)
	}

	next := scanDelimitedIdentifierBody(s)
	if next != nil {
		s.restore(chck)
		return next
	}

	if !s.accept(Quote) {
		s.restore(chck)
		return errorExpected(Quote)
	}
	return nil
}

func scanDelimitedIdentifierBody(s *Scanner) state {
	chck := s.checkpoint()
	atLeastOne := false

	for {
		next := scanDelimitedIdentifierPart(s)
		if next != nil {
			break
		}
		atLeastOne = true
	}

	if atLeastOne {
		return nil
	}

	s.restore(chck)
	return errorExpected("at least one DelimitedIdentifierPart")
}

func scanDelimitedIdentifierPart(s *Scanner) state {
	if next := scanDoubleQuoteSymbol(s); next == nil {
		return nil
	}

	// TODO: nondoublequote character, see
	// https://github.com/tomarrell/lbadd/projects/4#card-31143787
	if s.accept(SimpleLatinCharacter) {
		return nil
	}

	return errorExpected("non DoubleQuote character (current implementation: SimpleLatinCharacter) or \"\"")
}

func scanDoubleQuoteSymbol(s *Scanner) state {
	chck := s.checkpoint()
	if s.accept(Quote) {
		if s.accept(Quote) {
			return nil
		}
	}
	s.restore(chck)
	return errorExpected(Quote)
}

func scanUnicodeDelimitedIdentifier(s *Scanner) state {
	chck := s.checkpoint()

	if !s.acceptString("U") {
		s.restore(chck)
		return errorExpected("U")
	}

	if !s.accept(Ampersand) {
		s.restore(chck)
		return errorExpected(Ampersand)
	}

	if !s.accept(DoubleQuote) {
		s.restore(chck)
		return errorExpected(DoubleQuote)
	}

	if next := scanUnicodeDelimiterBody(s); next != nil {
		s.restore(chck)
		return next
	}

	if !s.accept(DoubleQuote) {
		s.restore(chck)
		return errorExpected(DoubleQuote)
	}

	if next := scanUnicodeEscapeSpecifier(s); next != nil {
		s.restore(chck)
		return next
	}

	return nil
}
