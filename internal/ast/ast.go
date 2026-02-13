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

// ParticipantKind classifies sequence diagram participant types.
type ParticipantKind int

const (
	ParticipantDefault ParticipantKind = iota
	ParticipantActor
	ParticipantBoundary
	ParticipantControl
	ParticipantEntity
	ParticipantDatabase
	ParticipantCollections
	ParticipantQueue
)

// FragmentKind classifies combined fragment types.
type FragmentKind int

const (
	FragmentAlt   FragmentKind = iota // alt/else
	FragmentLoop                      // loop
	FragmentPar                       // par
	FragmentGroup                     // group
)

// NotePosition indicates where a note is placed.
type NotePosition int

const (
	NoteLeft NotePosition = iota
	NoteRight
	NoteOver
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

// Member is the interface for class/interface/enum members.
type Member interface {
	Node
	memberNode()
}

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

// Participant represents a sequence diagram participant declaration.
type Participant struct {
	Pos   lexer.Pos
	Name  string
	Alias string
	Kind  ParticipantKind
}

func (p *Participant) Position() lexer.Pos { return p.Pos }
func (p *Participant) stmtNode()           {}

// Message represents a sequence diagram message between participants.
type Message struct {
	Pos    lexer.Pos
	From   string
	To     string
	Label  string
	Arrow  string // raw arrow literal
	Dashed bool
}

func (m *Message) Position() lexer.Pos { return m.Pos }
func (m *Message) stmtNode()           {}

// Fragment represents a combined fragment (alt, loop, par, group).
type Fragment struct {
	Pos        lexer.Pos
	Kind       FragmentKind
	Condition  string
	Statements []Statement
	ElseParts  []ElsePart
}

func (f *Fragment) Position() lexer.Pos { return f.Pos }
func (f *Fragment) stmtNode()           {}

// ElsePart represents an else clause within a fragment.
type ElsePart struct {
	Pos        lexer.Pos
	Condition  string
	Statements []Statement
}

func (e *ElsePart) Position() lexer.Pos { return e.Pos }

// Note represents a note attached to an element or floating.
type Note struct {
	Pos       lexer.Pos
	Placement NotePosition
	Target    string // element the note is attached to
	Text      string
}

func (n *Note) Position() lexer.Pos { return n.Pos }
func (n *Note) stmtNode()           {}

// Skinparam represents a skinparam directive.
type Skinparam struct {
	Pos   lexer.Pos
	Name  string
	Value string
}

func (s *Skinparam) Position() lexer.Pos { return s.Pos }
func (s *Skinparam) stmtNode()           {}

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

// HideShow represents a hide or show directive.
type HideShow struct {
	Pos    lexer.Pos
	IsHide bool // true for hide, false for show
	Target string
}

func (h *HideShow) Position() lexer.Pos { return h.Pos }
func (h *HideShow) stmtNode()           {}

// Activate represents an activate/deactivate statement.
type Activate struct {
	Pos        lexer.Pos
	Target     string
	Deactivate bool
}

func (a *Activate) Position() lexer.Pos { return a.Pos }
func (a *Activate) stmtNode()           {}

// Return represents a return message in a sequence diagram.
type Return struct {
	Pos   lexer.Pos
	Label string
}

func (r *Return) Position() lexer.Pos { return r.Pos }
func (r *Return) stmtNode()           {}
