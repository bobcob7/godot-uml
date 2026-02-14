package ast

import "github.com/bobcob7/godot-uml/internal/lexer"

// NotePosition indicates where a note is placed.
type NotePosition int

const (
	NoteLeft NotePosition = iota
	NoteRight
	NoteOver
)

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

// HideShow represents a hide or show directive.
type HideShow struct {
	Pos    lexer.Pos
	IsHide bool // true for hide, false for show
	Target string
}

func (h *HideShow) Position() lexer.Pos { return h.Pos }
func (h *HideShow) stmtNode()           {}
