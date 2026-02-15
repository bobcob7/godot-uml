package ast

import "github.com/bobcob7/go-uml/internal/lexer"

// RelationshipType classifies the kind of relationship between elements.
type RelationshipType int

const (
	RelAssociation RelationshipType = iota // --
	RelDependency                          // ..>
	RelInheritance                         // --|>
	RelRealization                         // ..|>
	RelComposition                         // --*
	RelAggregation                         // --o
)

// ArrowDirection indicates the directionality of a relationship arrow.
type ArrowDirection int

const (
	ArrowLeft  ArrowDirection = iota // <--
	ArrowRight                       // -->
	ArrowBoth                        // <-->
	ArrowNone                        // --
)

// Member is the interface for class/interface/enum members.
type Member interface {
	Node
	memberNode()
}

// ClassDef represents a class definition.
type ClassDef struct {
	Pos        lexer.Pos
	Name       string
	Alias      string
	Abstract   bool
	Members    []Member
	Stereotype string
}

func (c *ClassDef) Position() lexer.Pos { return c.Pos }
func (c *ClassDef) stmtNode()           {}

// InterfaceDef represents an interface definition.
type InterfaceDef struct {
	Pos        lexer.Pos
	Name       string
	Alias      string
	Members    []Member
	Stereotype string
}

func (i *InterfaceDef) Position() lexer.Pos { return i.Pos }
func (i *InterfaceDef) stmtNode()           {}

// EnumDef represents an enum definition.
type EnumDef struct {
	Pos        lexer.Pos
	Name       string
	Alias      string
	Values     []string
	Members    []Member
	Stereotype string
}

func (e *EnumDef) Position() lexer.Pos { return e.Pos }
func (e *EnumDef) stmtNode()           {}

// Field represents a class field/attribute.
type Field struct {
	Pos        lexer.Pos
	Name       string
	Type       string
	Visibility Visibility
	Modifier   Modifier
}

func (f *Field) Position() lexer.Pos { return f.Pos }
func (f *Field) memberNode()         {}

// Method represents a class method/operation.
type Method struct {
	Pos        lexer.Pos
	Name       string
	Params     string // raw parameter text
	ReturnType string
	Visibility Visibility
	Modifier   Modifier
}

func (m *Method) Position() lexer.Pos { return m.Pos }
func (m *Method) memberNode()         {}

// Relationship represents a connection between two elements.
type Relationship struct {
	Pos       lexer.Pos
	Left      string
	Right     string
	Type      RelationshipType
	Direction ArrowDirection
	Label     string
	LeftCard  string // left cardinality
	RightCard string // right cardinality
	Arrow     string // raw arrow literal
}

func (r *Relationship) Position() lexer.Pos { return r.Pos }
func (r *Relationship) stmtNode()           {}

// Package represents a package or namespace grouping.
type Package struct {
	Pos         lexer.Pos
	Name        string
	Alias       string
	Statements  []Statement
	IsNamespace bool
}

func (p *Package) Position() lexer.Pos { return p.Pos }
func (p *Package) stmtNode()           {}
