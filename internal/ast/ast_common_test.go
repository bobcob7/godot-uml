package ast_test

import (
	"testing"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/bobcob7/go-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func TestNoteStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 2, Column: 1}
		n := &ast.Note{Pos: pos, Placement: ast.NoteLeft, Target: "Foo", Text: "hello"}
		var s ast.Statement = n
		assert.Equal(t, pos, s.Position())
	})
}

func TestSkinparamStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 3, Column: 1}
		sp := &ast.Skinparam{Pos: pos, Name: "backgroundColor", Value: "#FFF"}
		var s ast.Statement = sp
		assert.Equal(t, pos, s.Position())
	})
}

func TestHideShowStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 4, Column: 1}
		hs := &ast.HideShow{Pos: pos, IsHide: true, Target: "members"}
		var s ast.Statement = hs
		assert.Equal(t, pos, s.Position())
	})
}

func TestNotePositionConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.NotePosition(0), ast.NoteLeft)
		assert.Equal(t, ast.NotePosition(1), ast.NoteRight)
		assert.Equal(t, ast.NotePosition(2), ast.NoteOver)
	})
}
