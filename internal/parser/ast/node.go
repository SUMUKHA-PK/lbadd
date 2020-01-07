package ast

type Node interface {
	Positioner
	Lengther
	Typer
	Valuer
}

type Lengther interface {
	Length() int
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

var _ Node = (*nodeInfo)(nil) // ensure nodeInfo implements Node

type nodeInfo struct {
	length int
	line   int
	col    int
	offset int
	typ    NodeType
	value  string
}

func (i nodeInfo) Length() int {
	return i.length
}

func (i nodeInfo) Line() int {
	return i.line
}

func (i nodeInfo) Col() int {
	return i.col
}

func (i nodeInfo) Offset() int {
	return i.offset
}

func (i nodeInfo) Type() NodeType {
	return i.typ
}

func (i nodeInfo) Value() string {
	return i.value
}

func (i nodeInfo) HasValue() bool {
	return i.value != ""
}
