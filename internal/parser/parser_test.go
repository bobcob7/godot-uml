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
		diagram, errs := Parse("@startuml\n@enduml")
		require.Empty(t, errs)
		assert.NotNil(t, diagram)
		assert.Empty(t, diagram.Statements)
		assert.Equal(t, lexer.Pos{Line: 1, Column: 1}, diagram.Pos)
	})
	t.Run("EmptyDiagramWithName", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml MyDiagram\n@enduml")
		require.Empty(t, errs)
		assert.Equal(t, "MyDiagram", diagram.Name)
	})
	t.Run("EmptyInput", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("")
		assert.NotNil(t, diagram)
		assert.Empty(t, errs)
	})
	t.Run("OnlyComments", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\n' comment\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		c, ok := diagram.Statements[0].(*ast.Comment)
		require.True(t, ok)
		assert.Equal(t, "' comment", c.Text)
	})
	t.Run("TitleDirective", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ntitle My Title\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
	})
	t.Run("SkinparamDirective", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nskinparam backgroundColor #FFF\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		sp, ok := diagram.Statements[0].(*ast.Skinparam)
		require.True(t, ok)
		assert.Equal(t, "backgroundColor", sp.Name)
	})
	t.Run("CaseInsensitiveStartUML", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@StartUml\n@EndUml")
		require.Empty(t, errs)
		assert.NotNil(t, diagram)
	})
}

func TestParseErrors(t *testing.T) {
	t.Parallel()
	t.Run("MissingStartUML", func(t *testing.T) {
		t.Parallel()
		_, errs := Parse("title Hello\n@enduml")
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Message, "expected @startuml")
	})
	t.Run("MissingEndUML", func(t *testing.T) {
		t.Parallel()
		_, errs := Parse("@startuml\ntitle Hello")
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[len(errs)-1].Message, "expected @enduml")
	})
	t.Run("ErrorRecovery", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\n$invalid\ntitle Valid\n@enduml")
		require.NotEmpty(t, errs)
		require.Len(t, diagram.Statements, 1)
	})
	t.Run("MultipleErrors", func(t *testing.T) {
		t.Parallel()
		_, errs := Parse("@startuml\n$one\n$two\n@enduml")
		assert.GreaterOrEqual(t, len(errs), 2)
	})
	t.Run("ErrorPositions", func(t *testing.T) {
		t.Parallel()
		_, errs := Parse("@startuml\n$bad\n@enduml")
		require.NotEmpty(t, errs)
		assert.Equal(t, 2, errs[0].Pos.Line)
	})
	t.Run("MissingClosingBrace", func(t *testing.T) {
		t.Parallel()
		_, errs := Parse("@startuml\nclass Foo {\n+name : String\n@enduml")
		require.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Message, "expected closing }")
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
	})
}
