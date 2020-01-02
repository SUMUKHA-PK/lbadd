package parser

import "github.com/tomarrell/lbadd/internal/parser/ast"

func Parse(sql string) ast.Query {
	return New().Parse(sql)
}

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(sql string) ast.Query {
	panic("TODO")
}
