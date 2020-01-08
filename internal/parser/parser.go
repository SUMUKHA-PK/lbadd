package parser

import (
	"fmt"

	"github.com/tomarrell/lbadd/internal/parser/ast"
)

const (
	ErrParserError          = sentinel("Parser error")
	ErrScannerError         = sentinel("Scanner error")
	ErrTimeout              = sentinel("Operation timed out")
	ErrUnexpectedToken      = sentinel("Unexpected token")
	ErrUnrecoverable        = sentinel("Unrecoverable error")
	ErrUnsupportedConstruct = sentinel("Construct not supported")
)

func Parse(sql string) (*ast.Query, error) {
	return New().Parse(sql)
}

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) (*ast.Query, error) {
	iso := newIsolate(sql)

	query, err := iso.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return query, nil
}
