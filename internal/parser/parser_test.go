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

func TestSeqModeAutoDetection(t *testing.T) {
	t.Parallel()
	t.Run("ParticipantTriggersSeqMode", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		_, ok := diagram.Statements[2].(*ast.Message)
		assert.True(t, ok, "ident->arrow->ident should produce Message in seqMode")
	})
	t.Run("ActorTriggersSeqMode", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nactor Alice\nactor Bob\nAlice -> Bob : hi\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		_, ok := diagram.Statements[2].(*ast.Message)
		assert.True(t, ok, "ident->arrow->ident should produce Message after actor declaration")
	})
	t.Run("WithoutSeqKeywordProducesRelationship", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nFoo --> Bar\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		_, ok := diagram.Statements[0].(*ast.Relationship)
		assert.True(t, ok, "ident->arrow->ident without seq keyword should produce Relationship")
	})
	t.Run("ActivateTriggersSeqMode", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nactivate Alice\nAlice -> Bob : msg\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		_, ok := diagram.Statements[1].(*ast.Message)
		assert.True(t, ok, "activate should trigger seqMode")
	})
	t.Run("AltTriggersSeqMode", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant A\nparticipant B\nalt test\nA -> B : msg\nend\nA -> B : after\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		var msgs int
		for _, stmt := range diagram.Statements {
			if _, ok := stmt.(*ast.Message); ok {
				msgs++
			}
		}
		assert.GreaterOrEqual(t, msgs, 1, "messages after fragment should still be parsed as Messages")
	})
}

func TestParseNoteShared(t *testing.T) {
	t.Parallel()
	t.Run("NoteLeftInClassContext", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\nnote left of Foo : a class note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteLeft, n.Placement)
		assert.Equal(t, "Foo", n.Target)
		assert.Equal(t, "a class note", n.Text)
	})
	t.Run("NoteRightInSequenceContext", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote right of Alice : a seq note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteRight, n.Placement)
		assert.Equal(t, "Alice", n.Target)
		assert.Equal(t, "a seq note", n.Text)
	})
	t.Run("NoteOverWithCommaTargets", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nnote over Alice,Bob : shared note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		n, ok := diagram.Statements[2].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteOver, n.Placement)
		assert.Equal(t, "Alice,Bob", n.Target)
		assert.Equal(t, "shared note", n.Text)
	})
	t.Run("MultiLineNoteInClassContext", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\nnote left of Foo\nLine one\nLine two\nend note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Contains(t, n.Text, "Line one")
		assert.Contains(t, n.Text, "Line two")
	})
	t.Run("MultiLineNoteInSequenceContext", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote right of Alice\nFirst line\nSecond line\nend note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Contains(t, n.Text, "First line")
		assert.Contains(t, n.Text, "Second line")
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
