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

// 7.2 row value expression
type (
	RowValueExpression struct {
		Node

		RowValueSpecialCase         *RowValueSpecialCase
		ExplicitRowValueConstructor *ExplicitRowValueConstructor
	}

	TableRowValueExpression struct {
		Node

		RowValueSpecialCase *RowValueSpecialCase
		RowValueConstructor *RowValueConstructor
	}

	ContextuallyTypedRowValueExpression struct {
		Node

		RowValueSpecialCase                 *RowValueSpecialCase
		ContextuallyTypedRowValueExpression *ContextuallyTypedRowValueExpression
	}

	RowValuePredicand struct {
		Node

		RowValueSpecialCase          *RowValueSpecialCase
		RowValueConstructorPredicand *RowValueConstructorPredicand
	}

	RowValueSpecialCase struct {
		Node

		NonparenthesizedValueExpressionPrimary *NonparenthesizedValueExpressionPrimary
	}
)

// 7.3 table value constructor
type (
	TableValueConstructor struct {
		Node

		Values                 token.Token
		RowValueExpressionList *RowValueExpressionList
	}

	RowValueExpressionList struct {
		Node

		TableRowValueExpression []*TableRowValueExpression
	}

	ContextuallyTypedTableValueConstructor struct {
		Node

		Values                                  token.Token
		ContextuallyTypedRowValueExpressionList *ContextuallyTypedRowValueExpressionList
	}

	ContextuallyTypedRowValueExpressionList struct {
		Node

		ContextuallyTypedRowValueExpression []*ContextuallyTypedRowValueExpression
	}
)

// 7.4 table expression
type (
	TableExpression struct {
		Node

		FromClause    *FromClause
		WhereClause   *WhereClause
		GroupByClause *GroupByClause
		HavingClause  *HavingClause
		WindowClause  *WindowClause
	}
)

// 7.5 from clause
type (
	FromClause struct {
		Node

		From               token.Token
		TableReferenceList *TableReferenceList
	}

	TableReferenceList struct {
		Node

		TableReference []*TableReference
	}
)

// 7.6 table reference
type (
	TableReference struct {
		Node

		TableFactor *TableFactor
		JoinedTable *JoinedTable
	}

	TableFactor struct {
		Node

		TablePrimary *TablePrimary
		SampleClause *SampleClause
	}

	SampleClause struct {
		Node

		Tablesample      token.Token
		SampleMethod     *SampleMethod
		LeftParen        token.Token
		SamplePercentage *SamplePercentage
		RightParen       token.Token
		RepeatableClause *RepeatableClause
	}

	SampleMethod struct {
		Node

		Bernoulli token.Token
		System    token.Token
	}

	RepeatableClause struct {
		Node

		Repeatable     token.Token
		LeftParen      token.Token
		RepeatArgument *RepeatArgument
		RightParen     token.Token
	}

	SamplePercentage struct {
		Node

		NumericValueExpression *NumericValueExpression
	}

	RepeatArgument struct {
		Node

		NumericValueExpression *NumericValueExpression
	}

	TablePrimary struct {
		Node

		TableOrQueryName                   *TableOrQueryName
		QuerySystemTimePeriodSpecification *QuerySystemTimePeriodSpecification
		DerivedTable                       *DerivedTable
		LateralDerivedTable                *LateralDerivedTable
		CollectionDerivedTable             *CollectionDerivedTable
		TableFunctionDerivedTable          *TableFunctionDerivedTable
		PTFDerivedTable                    *PTFDerivedTable
		OnlySpec                           *OnlySpec
		DataChangeDeltaTable               *DataChangeDeltaTable
		JSONTable                          *JSONTable
		JSONTablePrimitive                 *JSONTablePrimitive
		CorrelationOrRecognition           *CorrelationOrRecognition
		CorrelationName                    *CorrelationName
		ParenthesizedJoinedTable           *ParenthesizedJoinedTable
	}

	CorrelationOrRecognition struct {
		Node

		As                                 token.Token
		CorrelationName                    *CorrelationName
		ParenthesizedDerivedColumnList     *ParenthesizedDerivedColumnList
		RowPatternRecognitionClauseAndName *RowPatternRecognitionClauseAndName
	}

	QuerySystemTimePeriodSpecification struct {
		Node

		For                   token.Token
		SystemTime            token.Token
		As                    token.Token
		Of                    token.Token
		Between               token.Token
		AsymmetricOrSymmetric token.Token
		PointInTime1          *PointInTime
		PointInTime2          *PointInTime
		From                  token.Token
		To                    token.Token
	}

	PointInTime struct {
		Node

		DatetimeValueExpression *DatetimeValueExpression
	}

	OnlySpec struct {
		Node

		Only             token.Token
		LeftParen        token.Token
		TableOrQueryName *TableOrQueryName
		RightParen       token.Token
	}

	LateralDerivedTable struct {
		Node

		Lateral       token.Token
		TableSubquery *TableSubquery
	}

	CollectionDerivedTable struct {
		Node

		Unnest                    token.Token
		LeftParen                 token.Token
		CollectionValueExpression []*CollectionValueExpression
		RightParen                token.Token
		With                      token.Token
		Ordinality                token.Token
	}

	TableFunctionDerivedTable struct {
		Node

		Table                     token.Token
		LeftParen                 token.Token
		CollectionValueExpression *CollectionValueExpression
		RightParen                token.Token
	}

	DerivedTable struct {
		Node

		TableSubquery *TableSubquery
	}

	PTFDerivedTable struct {
		Node

		Table             token.Token
		LeftParen         token.Token
		RoutineInvocation *RoutineInvocation
		RightParen        token.Token
	}

	TableOrQueryName struct {
		Node

		TableName           *TableName
		TransitionTableName *TransitionTableName
		QueryName           *QueryName
	}
)
