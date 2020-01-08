package ast

import "github.com/tomarrell/lbadd/internal/parser/scanner/token"

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
	line   int
	col    int
	offset int
	length int
	typ    NodeType
	value  string
}

func NewNodeInfo(line, col, offset, length int, typ NodeType, value string) nodeInfo {
	return nodeInfo{
		line:   line,
		col:    col,
		offset: offset,
		length: length,
		typ:    typ,
		value:  value,
	}
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

var _ Node = (*tokenNode)(nil) // ensure nodeInfo implements Node

type tokenNode struct {
	token.Token
	typ NodeType
}

func NewTokenNode(t token.Token, typ NodeType) tokenNode {
	return tokenNode{
		Token: t,
		typ:   typ,
	}
}

func (n tokenNode) Type() NodeType {
	return n.typ
}

func (n tokenNode) HasValue() bool {
	return n.Value() != ""
}
