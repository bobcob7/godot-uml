package parser

import (
	"testing"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_EmptyDiagram(t *testing.T) {
	t.Parallel()
	input := "@startuml\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	assert.NotNil(t, diagram)
	assert.Empty(t, diagram.Statements)
	assert.Equal(t, lexer.Pos{Line: 1, Column: 1}, diagram.Pos)
}

func TestParse_EmptyDiagramWithName(t *testing.T) {
	t.Parallel()
	input := "@startuml MyDiagram\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	assert.Equal(t, "MyDiagram", diagram.Name)
	assert.Empty(t, diagram.Statements)
}

func TestParse_EmptyInput(t *testing.T) {
	t.Parallel()
	diagram, errs := Parse("")
	assert.NotNil(t, diagram)
	assert.Empty(t, diagram.Statements)
	assert.Empty(t, errs)
}

func TestParse_OnlyComments(t *testing.T) {
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
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()
	input := "@startuml\n/' block comment '/\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	c, ok := diagram.Statements[0].(*ast.Comment)
	require.True(t, ok)
	assert.Equal(t, "/' block comment '/", c.Text)
}

func TestParse_TitleDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\ntitle My Diagram Title\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	c, ok := diagram.Statements[0].(*ast.Comment)
	require.True(t, ok)
	assert.Equal(t, "title My Diagram Title", c.Text)
}

func TestParse_HeaderDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\nheader Page Header\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	c, ok := diagram.Statements[0].(*ast.Comment)
	require.True(t, ok)
	assert.Equal(t, "header Page Header", c.Text)
}

func TestParse_FooterDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\nfooter Page Footer\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	c, ok := diagram.Statements[0].(*ast.Comment)
	require.True(t, ok)
	assert.Equal(t, "footer Page Footer", c.Text)
}

func TestParse_SkinparamDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\nskinparam backgroundColor #FEFECE\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	sp, ok := diagram.Statements[0].(*ast.Skinparam)
	require.True(t, ok)
	assert.Equal(t, "backgroundColor", sp.Name)
	assert.Equal(t, "# FEFECE", sp.Value)
}

func TestParse_HideDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\nhide empty members\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	hs, ok := diagram.Statements[0].(*ast.HideShow)
	require.True(t, ok)
	assert.True(t, hs.IsHide)
	assert.Contains(t, hs.Target, "empty")
}

func TestParse_ShowDirective(t *testing.T) {
	t.Parallel()
	input := "@startuml\nshow methods\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	hs, ok := diagram.Statements[0].(*ast.HideShow)
	require.True(t, ok)
	assert.False(t, hs.IsHide)
}

func TestParse_MissingStartUML(t *testing.T) {
	t.Parallel()
	input := "title Hello\n@enduml"
	_, errs := Parse(input)
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Message, "expected @startuml")
}

func TestParse_MissingEndUML(t *testing.T) {
	t.Parallel()
	input := "@startuml\ntitle Hello"
	_, errs := Parse(input)
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[len(errs)-1].Message, "expected @enduml")
}

func TestParse_ErrorRecovery(t *testing.T) {
	t.Parallel()
	input := "@startuml\n$invalid\ntitle Valid Title\n@enduml"
	diagram, errs := Parse(input)
	require.NotEmpty(t, errs, "should report error for invalid token")
	require.Len(t, diagram.Statements, 1, "should recover and parse title")
	c, ok := diagram.Statements[0].(*ast.Comment)
	require.True(t, ok)
	assert.Equal(t, "title Valid Title", c.Text)
}

func TestParse_MultipleErrors(t *testing.T) {
	t.Parallel()
	input := "@startuml\n$one\n$two\ntitle OK\n@enduml"
	diagram, errs := Parse(input)
	assert.GreaterOrEqual(t, len(errs), 2, "should report multiple errors")
	require.Len(t, diagram.Statements, 1, "should still parse valid statement")
}

func TestParse_ErrorPositions(t *testing.T) {
	t.Parallel()
	input := "@startuml\n$bad\n@enduml"
	_, errs := Parse(input)
	require.NotEmpty(t, errs)
	assert.Equal(t, 2, errs[0].Pos.Line, "error should be on line 2")
}

func TestParse_ErrorMessage(t *testing.T) {
	t.Parallel()
	input := "@startuml\n$bad\n@enduml"
	_, errs := Parse(input)
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "2:")
	assert.Contains(t, errs[0].Error(), "unexpected")
}

func TestParse_MultipleDirectives(t *testing.T) {
	t.Parallel()
	input := "@startuml\ntitle My Title\n' a comment\nskinparam fontSize 14\nfooter Done\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 4)
}

func TestParse_WhitespaceLines(t *testing.T) {
	t.Parallel()
	input := "@startuml\n\n\ntitle Hello\n\n\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
}

func TestParse_CommentsAndDirectivesMixed(t *testing.T) {
	t.Parallel()
	input := "@startuml\n' comment 1\ntitle My Title\n' comment 2\nfooter Done\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 4)
}

func TestParse_SkinparamNoValue(t *testing.T) {
	t.Parallel()
	input := "@startuml\nskinparam shadowing\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	require.Len(t, diagram.Statements, 1)
	sp, ok := diagram.Statements[0].(*ast.Skinparam)
	require.True(t, ok)
	assert.Equal(t, "shadowing", sp.Name)
	assert.Empty(t, sp.Value)
}

func TestParse_DiagramNameQuoted(t *testing.T) {
	t.Parallel()
	input := "@startuml \"my diagram\"\n@enduml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	assert.Equal(t, "\"my diagram\"", diagram.Name)
}

func TestNew_AcceptsTokenSlice(t *testing.T) {
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
}

func TestParse_CaseInsensitiveStartUML(t *testing.T) {
	t.Parallel()
	input := "@StartUml\n@EndUml"
	diagram, errs := Parse(input)
	require.Empty(t, errs)
	assert.NotNil(t, diagram)
}
