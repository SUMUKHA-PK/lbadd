package ast

import "github.com/tomarrell/lbadd/internal/parser/scanner/token"

type (
	ScalarSubquery struct {
		Node

		Subquery *Subquery
	}

	RowSubquery struct {
		Node

		Subquery *Subquery
	}

	TableSubquery struct {
		Node

		Subquery *Subquery
	}

	Subquery struct {
		Node

		LeftParen       token.Token
		QueryExpression *QueryExpression
		RightParen      token.Token
	}
)
