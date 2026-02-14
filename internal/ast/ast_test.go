package ast_test

import (
	"testing"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func TestDiagram(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsNode", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 3, Column: 7}
		d := &ast.Diagram{Pos: pos}
		var n ast.Node = d
		assert.Equal(t, pos, n.Position())
	})
}

func TestComment(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsNode", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 5, Column: 1}
		c := &ast.Comment{Pos: pos, Text: "a comment"}
		var n ast.Node = c
		assert.Equal(t, pos, n.Position())
	})
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		c := &ast.Comment{Pos: lexer.Pos{Line: 1, Column: 1}}
		var s ast.Statement = c
		assert.Equal(t, lexer.Pos{Line: 1, Column: 1}, s.Position())
	})
}
