package ast

import "github.com/tomarrell/lbadd/internal/parser/scanner/token"

// 7 Query expressions

// 7.1 row value constructor
type (
	RowValueConstructor struct {
		Node

		CommonValueExpression       *CommonValueExpression
		BooleanValueExpression      *BooleanValueExpression
		ExplicitRowValueConstructor *ExplicitRowValueConstructor
	}

	ExplicitRowValueConstructor struct {
		Node

		LeftParen                      token.Token
		RowValueConstructorElement     *RowValueConstructorElement
		Comma                          token.Token
		RowValueConstructorElementList *RowValueConstructorElementList
		RightParen                     token.Token
		Row                            token.Token
		RowSubquery                    *RowSubquery
	}

	RowValueConstructorElementList struct {
		Node

		RowValueConstructorElement []*RowValueConstructorElement
	}

	RowValueConstructorElement struct {
		Node

		ValueExpression *ValueExpression
	}

	ContextuallyTypedRowValueConstructorElementList struct {
		Node

		ContextuallyTypedRowValueConstructorElement []*ContextuallyTypedRowValueConstructorElement
	}

	ContextuallyTypedRowValueConstructorElement struct {
		Node

		ValueExpression                     *ValueExpression
		ContextuallyTypedValueSpecification *ContextuallyTypedValueSpecification
	}

	RowValueConstructorPredicand struct {
		Node

		CommonValueExpression       *CommonValueExpression
		BooleanPredicand            *BooleanPredicand
		ExplicitRowValueConstructor *ExplicitRowValueConstructor
	}
)
