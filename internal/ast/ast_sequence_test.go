package ast_test

import (
	"testing"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParticipantStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 2, Column: 1}
		p := &ast.Participant{Pos: pos, Name: "Alice", Kind: ast.ParticipantDefault}
		var s ast.Statement = p
		assert.Equal(t, pos, s.Position())
	})
}

func TestMessageStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 3, Column: 1}
		m := &ast.Message{Pos: pos, From: "Alice", To: "Bob"}
		var s ast.Statement = m
		assert.Equal(t, pos, s.Position())
	})
}

func TestFragmentStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 4, Column: 1}
		f := &ast.Fragment{Pos: pos, Kind: ast.FragmentAlt}
		var s ast.Statement = f
		assert.Equal(t, pos, s.Position())
	})
}

func TestActivateStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 5, Column: 1}
		a := &ast.Activate{Pos: pos, Target: "Bob"}
		var s ast.Statement = a
		assert.Equal(t, pos, s.Position())
	})
}

func TestReturnStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 6, Column: 1}
		r := &ast.Return{Pos: pos, Label: "ok"}
		var s ast.Statement = r
		assert.Equal(t, pos, s.Position())
	})
}

func TestAutonumberStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 7, Column: 1}
		a := &ast.Autonumber{Pos: pos}
		var s ast.Statement = a
		assert.Equal(t, pos, s.Position())
	})
}

func TestDividerStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 8, Column: 1}
		d := &ast.Divider{Pos: pos, Text: "Init"}
		var s ast.Statement = d
		assert.Equal(t, pos, s.Position())
	})
}

func TestDelayStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 9, Column: 1}
		d := &ast.Delay{Pos: pos}
		var s ast.Statement = d
		assert.Equal(t, pos, s.Position())
	})
}

func TestElsePartNode(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsNode", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 10, Column: 1}
		ep := ast.ElsePart{Pos: pos, Condition: "failure"}
		var n ast.Node = &ep
		assert.Equal(t, pos, n.Position())
	})
}

func TestParticipantKindConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.ParticipantKind(0), ast.ParticipantDefault)
		assert.Equal(t, ast.ParticipantKind(1), ast.ParticipantActor)
		assert.Equal(t, ast.ParticipantKind(2), ast.ParticipantBoundary)
		assert.Equal(t, ast.ParticipantKind(3), ast.ParticipantControl)
		assert.Equal(t, ast.ParticipantKind(4), ast.ParticipantEntity)
		assert.Equal(t, ast.ParticipantKind(5), ast.ParticipantDatabase)
		assert.Equal(t, ast.ParticipantKind(6), ast.ParticipantCollections)
		assert.Equal(t, ast.ParticipantKind(7), ast.ParticipantQueue)
	})
}

func TestFragmentKindConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.FragmentKind(0), ast.FragmentAlt)
		assert.Equal(t, ast.FragmentKind(1), ast.FragmentLoop)
		assert.Equal(t, ast.FragmentKind(2), ast.FragmentPar)
		assert.Equal(t, ast.FragmentKind(3), ast.FragmentBreak)
		assert.Equal(t, ast.FragmentKind(4), ast.FragmentRef)
		assert.Equal(t, ast.FragmentKind(5), ast.FragmentGroup)
	})
}
