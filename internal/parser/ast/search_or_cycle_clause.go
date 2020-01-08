package ast

import "github.com/tomarrell/lbadd/internal/parser/scanner/token"

type (
	SearchOrCycleClause struct {
		Node

		SearchClause *SearchClause
		CycleClause  *CycleClause
	}

	SearchClause struct {
		Node

		Search               token.Token
		RecursiveSearchOrder *RecursiveSearchOrder
		Set                  token.Token
		SequenceColumn       *SequenceColumn
	}

	RecursiveSearchOrder struct {
		Node

		DepthOrBreadth token.Token
		First          token.Token
		By             token.Token
		ColumnNameList *ColumnNameList
	}

	SequenceColumn struct {
		Node

		ColumnName *ColumnName
	}

	CycleClause struct {
		Node

		Cycle             token.Token
		CycleColumnList   *CycleColumnList
		Set               token.Token
		CycleMarkColumn   *CycleMarkColumn
		To                token.Token
		CycleMarkValue    *CycleMarkValue
		Default           token.Token
		NonCycleMarkValue *NonCycleMarkValue
		Using             token.Token
		PathColumn        *PathColumn
	}

	CycleColumnList struct {
		Node

		CycleColumn []*CycleColumn
	}

	CycleColumn struct {
		Node

		ColumnName *ColumnName
	}

	CycleMarkColumn struct {
		Node

		ColumnName *ColumnName
	}

	PathColumn struct {
		Node

		ColumnName *ColumnName
	}

	CycleMarkValue struct {
		Node

		ValueExpression *ValueExpression
	}

	NonCycleMarkValue struct {
		Node

		ValueExpression *ValueExpression
	}
)
