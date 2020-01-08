package ast

import "github.com/tomarrell/lbadd/internal/parser/scanner/token"

type (
	WithClause struct {
		Node

		With      token.Token
		Recursive token.Token
		WithList  *WithList
	}

	WithList struct {
		WithListElement []*WithListElement
	}

	WithListElement struct {
		Node

		QueryName           *QueryName
		LeftParen           token.Token
		WithColumnList      *WithColumnList
		RightParen          token.Token
		As                  token.Token
		TableSubquery       *TableSubquery
		SearchOrCycleClause *SearchOrCycleClause
	}

	WithColumnList struct {
		Node

		ColumnNameList *ColumnNameList
	}
)
