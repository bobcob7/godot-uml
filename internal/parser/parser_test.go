package parser

import (
	"testing"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()
	t.Run("EmptyDiagram", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		assert.NotNil(t, diagram)
		assert.Empty(t, diagram.Statements)
		assert.Equal(t, lexer.Pos{Line: 1, Column: 1}, diagram.Pos)
	})
	t.Run("EmptyDiagramWithName", func(t *testing.T) {
		t.Parallel()
		input := "@startuml MyDiagram\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		assert.Equal(t, "MyDiagram", diagram.Name)
		assert.Empty(t, diagram.Statements)
	})
	t.Run("EmptyInput", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("")
		assert.NotNil(t, diagram)
		assert.Empty(t, diagram.Statements)
		assert.Empty(t, errs)
	})
	t.Run("OnlyComments", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n' this is a comment\n' another comment\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		c1, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "' this is a comment", c1.Text)
		c2, ok := diagram.Statements[1].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "' another comment", c2.Text)
	})
	t.Run("BlockComment", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n/' block comment '/\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "/' block comment '/", c.Text)
	})
	t.Run("TitleDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\ntitle My Diagram Title\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "title My Diagram Title", c.Text)
	})
	t.Run("HeaderDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nheader Page Header\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "header Page Header", c.Text)
	})
	t.Run("FooterDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nfooter Page Footer\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "footer Page Footer", c.Text)
	})
	t.Run("SkinparamDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nskinparam backgroundColor #FEFECE\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		sp, ok := diagram.Statements[0].(*ast.Skinparam)
		require.True(t, ok)
		assert.Equal(t, "backgroundColor", sp.Name)
		assert.Equal(t, "# FEFECE", sp.Value)
	})
	t.Run("HideDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nhide empty members\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		hs, ok := diagram.Statements[0].(*ast.HideShow)
		require.True(t, ok)
		assert.True(t, hs.IsHide)
		assert.Contains(t, hs.Target, "empty")
	})
	t.Run("ShowDirective", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nshow methods\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		hs, ok := diagram.Statements[0].(*ast.HideShow)
		require.True(t, ok)
		assert.False(t, hs.IsHide)
	})
	t.Run("MultipleDirectives", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\ntitle My Title\n' a comment\nskinparam fontSize 14\nfooter Done\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 4)
	})
	t.Run("WhitespaceLines", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n\n\ntitle Hello\n\n\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
	})
	t.Run("CommentsAndDirectivesMixed", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n' comment 1\ntitle My Title\n' comment 2\nfooter Done\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 4)
	})
	t.Run("SkinparamNoValue", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nskinparam shadowing\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		sp, ok := diagram.Statements[0].(*ast.Skinparam)
		require.True(t, ok)
		assert.Equal(t, "shadowing", sp.Name)
		assert.Empty(t, sp.Value)
	})
	t.Run("DiagramNameQuoted", func(t *testing.T) {
		t.Parallel()
		input := "@startuml \"my diagram\"\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		assert.Equal(t, "\"my diagram\"", diagram.Name)
	})
	t.Run("CaseInsensitiveStartUML", func(t *testing.T) {
		t.Parallel()
		input := "@StartUml\n@EndUml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		assert.NotNil(t, diagram)
	})
}

func TestParseErrors(t *testing.T) {
	t.Parallel()
	t.Run("MissingStartUML", func(t *testing.T) {
		t.Parallel()
		input := "title Hello\n@enduml"
		_, errs := Parse(input)
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Message, "expected @startuml")
	})
	t.Run("MissingEndUML", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\ntitle Hello"
		_, errs := Parse(input)
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[len(errs)-1].Message, "expected @enduml")
	})
	t.Run("ErrorRecovery", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n$invalid\ntitle Valid Title\n@enduml"
		diagram, errs := Parse(input)
		require.NotEmpty(t, errs, "should report error for invalid token")
		require.Len(t, diagram.Statements, 1, "should recover and parse title")
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "title Valid Title", c.Text)
	})
	t.Run("MultipleErrors", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n$one\n$two\ntitle OK\n@enduml"
		diagram, errs := Parse(input)
		assert.GreaterOrEqual(t, len(errs), 2, "should report multiple errors")
		require.Len(t, diagram.Statements, 1, "should still parse valid statement")
	})
	t.Run("ErrorPositions", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n$bad\n@enduml"
		_, errs := Parse(input)
		require.NotEmpty(t, errs)
		assert.Equal(t, 2, errs[0].Pos.Line, "error should be on line 2")
	})
	t.Run("ErrorMessage", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\n$bad\n@enduml"
		_, errs := Parse(input)
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Error(), "2:")
		assert.Contains(t, errs[0].Error(), "unexpected")
	})
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("AcceptsTokenSlice", func(t *testing.T) {
		t.Parallel()
		tokens := []lexer.Token{
			{Type: lexer.TokenStartUML, Literal: "@startuml", Pos: lexer.Pos{Line: 1, Column: 1}},
			{Type: lexer.TokenNewline, Literal: "\n", Pos: lexer.Pos{Line: 1, Column: 10}},
			{Type: lexer.TokenEndUML, Literal: "@enduml", Pos: lexer.Pos{Line: 2, Column: 1}},
			{Type: lexer.TokenEOF, Pos: lexer.Pos{Line: 2, Column: 8}},
		}
		p := New(tokens)
		diagram := p.parseDiagram()
		assert.Empty(t, p.Errors())
		assert.NotNil(t, diagram)
		assert.Empty(t, diagram.Statements)
	})
}
