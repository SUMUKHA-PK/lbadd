package ast

type Node interface {
	Positioner
	Typer
}

type ValueNode interface {
	Node
	Valuer
}

type Positioner interface {
	Line() int
	Col() int
	Offset() int
}

type Typer interface {
	Type() NodeType
}

type Valuer interface {
	Value() string
	HasValue() bool
}
