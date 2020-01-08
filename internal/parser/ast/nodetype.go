package ast

//go:generate stringer -type=NodeType
type NodeType uint16

const (
	Unknown NodeType = iota
	Root
)
