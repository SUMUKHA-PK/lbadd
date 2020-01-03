package ast

type Node interface {
	Positioner
	Typer

	// should Node also embed Valuer? maybe don't force the implementation but
	// leave it up to the node to implement Valuer. Conversion from AST to IR
	// can make it more verbose, but here, that also saves us having to keep
	// constant keywords in memory
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
