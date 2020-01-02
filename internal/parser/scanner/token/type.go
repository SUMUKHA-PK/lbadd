package token

//go:generate stringer -type=Type
type Type uint16

const (
	Unknown Type = iota
	// Error indicates that a syntax error has been detected by the lexical
	// analyzer (scanner) and that the error should be printed. The parser also
	// should consider resetting its state.
	Error
	// EOF indicates that the lexical analyzer (scanner) has reached the end of
	// its input. After receiving this token, the parser can close the token
	// stream, as there will not be any more tokens. He also can start building
	// (if not already done) the AST, as he know knows of all tokens.
	EOF

	// non-delimiter tokens

	RegularIdentifier
	KeyWord
	UnsignedNumericLiteral
	NationalCharacterStringLiteral
	BinaryStringLiteral
	LargeObjectLengthToken
	UnicodeDelimitedIdentifier
	UnicodeCharacterStringLiteral
	SQLLanguageIdentifier

	// delimiter tokens

	CharacterStringLiteral
	DateString
	TimeString
	TimestampString
	IntervalString
	DelimitedIdentifier
	SQLSpecialCharacter
	NotEqualsOperator
	GreaterThanOrEqualsOperator
	LessThanOrEqualsOperator
	ConcatenationOperator
	RightArrow
	LeftBracketTrigraph
	RightBracketTrigraph
	DoubleColon
	DoublePeriod
	NamedArgumentAssignmentToken
	LeftBraceMinus
	RightMinusBrace
)

func (t Type) IsNonDelimiterToken() bool {
	return t >= RegularIdentifier &&
		t <= SQLLanguageIdentifier
}

func (t Type) IsDelimiterToken() bool {
	return t >= CharacterStringLiteral &&
		t <= RightMinusBrace
}
