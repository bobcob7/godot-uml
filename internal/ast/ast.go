// Package ast defines the abstract syntax tree node types for PlantUML diagrams.
package ast

import "github.com/bobcob7/godot-uml/internal/lexer"

// Node is the interface implemented by all AST nodes.
type Node interface {
	// Position returns the source position of this node.
	Position() lexer.Pos
}

// Statement is the interface for top-level diagram statements.
type Statement interface {
	Node
	stmtNode()
}

// Visibility represents member access visibility.
type Visibility int

const (
	VisibilityNone      Visibility = iota
	VisibilityPublic               // +
	VisibilityPrivate              // -
	VisibilityProtected            // #
	VisibilityPackage              // ~
)

// Modifier represents a member modifier.
type Modifier int

const (
	ModifierNone   Modifier = iota
	ModifierStatic          // {static}
	ModifierField           // {field}
	ModifierMethod          // {method}
)

// Diagram is the root AST node representing a complete PlantUML diagram.
type Diagram struct {
	Pos        lexer.Pos
	Name       string // optional name after @startuml
	Title      string
	Header     string
	Footer     string
	Statements []Statement
}

func (d *Diagram) Position() lexer.Pos { return d.Pos }

// Comment represents a comment statement preserved in the AST.
type Comment struct {
	Pos  lexer.Pos
	Text string
}

func (c *Comment) Position() lexer.Pos { return c.Pos }
func (c *Comment) stmtNode()           {}
