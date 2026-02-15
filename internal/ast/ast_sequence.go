package ast

import "github.com/bobcob7/go-uml/internal/lexer"

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
	FragmentBreak                     // break
	FragmentRef                       // ref
	FragmentGroup                     // group
)

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

// Autonumber represents an autonumber directive in a sequence diagram.
type Autonumber struct {
	Pos   lexer.Pos
	Start string
}

func (a *Autonumber) Position() lexer.Pos { return a.Pos }
func (a *Autonumber) stmtNode()           {}

// Divider represents a divider (== text ==) in a sequence diagram.
type Divider struct {
	Pos  lexer.Pos
	Text string
}

func (d *Divider) Position() lexer.Pos { return d.Pos }
func (d *Divider) stmtNode()           {}

// Delay represents a delay (...) or (... text ...) in a sequence diagram.
type Delay struct {
	Pos  lexer.Pos
	Text string
}

func (d *Delay) Position() lexer.Pos { return d.Pos }
func (d *Delay) stmtNode()           {}
