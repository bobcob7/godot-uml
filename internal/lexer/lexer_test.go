package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextToken_DiagramDelimiters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantTok Token
	}{
		{
			name:    "startuml",
			input:   "@startuml",
			wantTok: Token{Type: TokenStartUML, Literal: "@startuml", Pos: Pos{1, 1}},
		},
		{
			name:    "enduml",
			input:   "@enduml",
			wantTok: Token{Type: TokenEndUML, Literal: "@enduml", Pos: Pos{1, 1}},
		},
		{
			name:    "startuml case insensitive",
			input:   "@StartUml",
			wantTok: Token{Type: TokenStartUML, Literal: "@StartUml", Pos: Pos{1, 1}},
		},
		{
			name:    "unknown directive is error",
			input:   "@foo",
			wantTok: Token{Type: TokenError, Literal: "@foo", Pos: Pos{1, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantTok, tok)
		})
	}
}

func TestNextToken_Keywords(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		wantType TokenType
	}{
		{"class", "class", TokenClass},
		{"interface", "interface", TokenInterface},
		{"enum", "enum", TokenEnum},
		{"abstract", "abstract", TokenAbstract},
		{"extends", "extends", TokenExtends},
		{"implements", "implements", TokenImplements},
		{"package", "package", TokenPackage},
		{"namespace", "namespace", TokenNamespace},
		{"as", "as", TokenAs},
		{"participant", "participant", TokenParticipant},
		{"actor", "actor", TokenActor},
		{"boundary", "boundary", TokenBoundary},
		{"control", "control", TokenControl},
		{"entity", "entity", TokenEntity},
		{"database", "database", TokenDatabase},
		{"collections", "collections", TokenCollections},
		{"queue", "queue", TokenQueue},
		{"activate", "activate", TokenActivate},
		{"deactivate", "deactivate", TokenDeactivate},
		{"return", "return", TokenReturn},
		{"alt", "alt", TokenAlt},
		{"else", "else", TokenElse},
		{"end", "end", TokenEnd},
		{"loop", "loop", TokenLoop},
		{"group", "group", TokenGroup},
		{"note", "note", TokenNote},
		{"of", "of", TokenOf},
		{"over", "over", TokenOver},
		{"left", "left", TokenLeft},
		{"right", "right", TokenRight},
		{"skinparam", "skinparam", TokenSkinparam},
		{"hide", "hide", TokenHide},
		{"show", "show", TokenShow},
		{"title", "title", TokenTitle},
		{"header", "header", TokenHeader},
		{"footer", "footer", TokenFooter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantType, tok.Type)
			assert.Equal(t, tt.input, tok.Literal)
		})
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantLit string
	}{
		{"simple", "MyClass", "MyClass"},
		{"with underscore", "my_var", "my_var"},
		{"with digits", "Class1", "Class1"},
		{"underscore prefix", "_private", "_private"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, TokenIdent, tok.Type)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestNextToken_Strings(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantTok Token
	}{
		{
			name:    "double quoted",
			input:   `"hello world"`,
			wantTok: Token{Type: TokenString, Literal: `"hello world"`, Pos: Pos{1, 1}},
		},
		{
			name:    "escaped quote",
			input:   `"say \"hi\""`,
			wantTok: Token{Type: TokenString, Literal: `"say \"hi\""`, Pos: Pos{1, 1}},
		},
		{
			name:    "unterminated string",
			input:   `"hello`,
			wantTok: Token{Type: TokenError, Literal: `"hello`, Pos: Pos{1, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantTok, tok)
		})
	}
}

func TestNextToken_Numbers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantLit string
	}{
		{"integer", "42", "42"},
		{"decimal", "3.14", "3.14"},
		{"zero", "0", "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, TokenNumber, tok.Type)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestNextToken_Punctuation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		wantType TokenType
		wantLit  string
	}{
		{"{", TokenLBrace, "{"},
		{"}", TokenRBrace, "}"},
		{"(", TokenLParen, "("},
		{")", TokenRParen, ")"},
		{"[", TokenLBracket, "["},
		{"]", TokenRBracket, "]"},
		{":", TokenColon, ":"},
		{",", TokenComma, ","},
		{".", TokenDot, "."},
		{"|", TokenPipe, "|"},
		{"=", TokenEquals, "="},
		{";", TokenSemicolon, ";"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantType, tok.Type)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestNextToken_VisibilityMarkers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		wantType TokenType
	}{
		{"+", TokenPlus},
		{"-", TokenMinus},
		{"#", TokenHash},
		{"~", TokenTilde},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantType, tok.Type)
			assert.Equal(t, tt.input, tok.Literal)
		})
	}
}

func TestNextToken_Modifiers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantLit  string
	}{
		{"static", "{static}", TokenStatic, "{static}"},
		{"field", "{field}", TokenField, "{field}"},
		{"method", "{method}", TokenMethod, "{method}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, tt.wantType, tok.Type)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestNextToken_Arrows(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantLit string
	}{
		{"solid right", "->", "->"},
		{"dashed right", "-->", "-->"},
		{"solid left", "<-", "<-"},
		{"dashed left", "<--", "<--"},
		{"inheritance", "<|--", "<|--"},
		{"composition", "*--", "*--"},
		{"aggregation", "o--", "o--"},
		{"solid right long", "--->", "--->"},
		{"dotted right", "..>", "..>"},
		{"long dotted right", "...>", "...>"},
		{"realization left", "<|..", "<|.."},
		{"composition right", "--*", "--*"},
		{"aggregation right", "--o", "--o"},
		{"bare solid", "--", "--"},
		{"bare dotted", "..", ".."},
		{"right inheritance", "--|>", "--|>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, TokenArrow, tok.Type, "input: %s", tt.input)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestNextToken_Comments(t *testing.T) {
	t.Parallel()
	t.Run("line comment", func(t *testing.T) {
		t.Parallel()
		l := New("' this is a comment")
		tok := l.NextToken()
		assert.Equal(t, TokenLineComment, tok.Type)
		assert.Equal(t, "' this is a comment", tok.Literal)
	})
	t.Run("block comment", func(t *testing.T) {
		t.Parallel()
		l := New("/' block comment '/")
		tok := l.NextToken()
		assert.Equal(t, TokenBlockComment, tok.Type)
		assert.Equal(t, "/' block comment '/", tok.Literal)
	})
	t.Run("unterminated block comment", func(t *testing.T) {
		t.Parallel()
		l := New("/' unterminated")
		tok := l.NextToken()
		assert.Equal(t, TokenError, tok.Type)
	})
}

func TestNextToken_Newlines(t *testing.T) {
	t.Parallel()
	l := New("a\nb")
	tok1 := l.NextToken()
	assert.Equal(t, TokenIdent, tok1.Type)
	assert.Equal(t, Pos{1, 1}, tok1.Pos)
	tok2 := l.NextToken()
	assert.Equal(t, TokenNewline, tok2.Type)
	tok3 := l.NextToken()
	assert.Equal(t, TokenIdent, tok3.Type)
	assert.Equal(t, "b", tok3.Literal)
	assert.Equal(t, Pos{2, 1}, tok3.Pos)
}

func TestTokenize_ClassDiagram(t *testing.T) {
	t.Parallel()
	input := `@startuml
class Animal {
  +name : String
  #age : int
  +speak() : void
}
class Dog extends Animal {
  +breed : String
}
Dog --|> Animal
@enduml`
	l := New(input)
	tokens := l.Tokenize()
	require.NotEmpty(t, tokens)
	assert.Equal(t, TokenStartUML, tokens[0].Type)
	last := tokens[len(tokens)-1]
	assert.Equal(t, TokenEOF, last.Type)
	// Find the @enduml.
	var foundEndUML bool
	for _, tok := range tokens {
		if tok.Type == TokenEndUML {
			foundEndUML = true
			break
		}
	}
	assert.True(t, foundEndUML, "should find @enduml token")
}

func TestTokenize_SequenceDiagram(t *testing.T) {
	t.Parallel()
	input := `@startuml
participant Alice
participant Bob
Alice -> Bob : Hello
Bob --> Alice : Hi
@enduml`
	l := New(input)
	tokens := l.Tokenize()
	require.NotEmpty(t, tokens)
	// Check for arrow tokens.
	var arrows []Token
	for _, tok := range tokens {
		if tok.Type == TokenArrow {
			arrows = append(arrows, tok)
		}
	}
	require.Len(t, arrows, 2)
	assert.Equal(t, "->", arrows[0].Literal)
	assert.Equal(t, "-->", arrows[1].Literal)
}

func TestNextToken_SourcePositions(t *testing.T) {
	t.Parallel()
	input := "class Foo {\n  +bar\n}"
	l := New(input)
	tok := l.NextToken() // class
	assert.Equal(t, Pos{1, 1}, tok.Pos)
	tok = l.NextToken() // Foo
	assert.Equal(t, Pos{1, 7}, tok.Pos)
	tok = l.NextToken() // {
	assert.Equal(t, Pos{1, 11}, tok.Pos)
	tok = l.NextToken() // \n
	assert.Equal(t, TokenNewline, tok.Type)
	tok = l.NextToken() // +
	assert.Equal(t, Pos{2, 3}, tok.Pos)
	tok = l.NextToken() // bar
	assert.Equal(t, Pos{2, 4}, tok.Pos)
	tok = l.NextToken() // \n
	assert.Equal(t, TokenNewline, tok.Type)
	tok = l.NextToken() // }
	assert.Equal(t, Pos{3, 1}, tok.Pos)
}

func TestNextToken_ErrorTokens(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
	}{
		{"unknown char", "$"},
		{"backtick", "`"},
		{"at with unknown", "@xyz"},
		{"unterminated string", `"hello`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, TokenError, tok.Type, "input: %s", tt.input)
		})
	}
}

func TestNextToken_EmptyInput(t *testing.T) {
	t.Parallel()
	l := New("")
	tok := l.NextToken()
	assert.Equal(t, TokenEOF, tok.Type)
}

func TestNextToken_WhitespaceOnly(t *testing.T) {
	t.Parallel()
	l := New("   \t  ")
	tok := l.NextToken()
	assert.Equal(t, TokenEOF, tok.Type)
}

func TestNextToken_SkipsSpacesAndTabs(t *testing.T) {
	t.Parallel()
	l := New("  class  Foo")
	tok1 := l.NextToken()
	assert.Equal(t, TokenClass, tok1.Type)
	tok2 := l.NextToken()
	assert.Equal(t, TokenIdent, tok2.Type)
	assert.Equal(t, "Foo", tok2.Literal)
}

func TestTokenize_Skinparam(t *testing.T) {
	t.Parallel()
	input := "skinparam backgroundColor #FEFECE"
	l := New(input)
	tokens := l.Tokenize()
	require.True(t, len(tokens) >= 3)
	assert.Equal(t, TokenSkinparam, tokens[0].Type)
	assert.Equal(t, TokenIdent, tokens[1].Type)
	assert.Equal(t, "backgroundColor", tokens[1].Literal)
}

func TestTokenize_MultipleArrowTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantLit string
	}{
		{"composition left", "*--", "*--"},
		{"aggregation left", "o--", "o--"},
		{"realization dotted", "<|..", "<|.."},
		{"dependency", "..>", "..>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := New(tt.input)
			tok := l.NextToken()
			assert.Equal(t, TokenArrow, tok.Type)
			assert.Equal(t, tt.wantLit, tok.Literal)
		})
	}
}

func TestTokenize_BraceModifiers(t *testing.T) {
	t.Parallel()
	input := "{static} name\n{field} x\n{method} foo()"
	l := New(input)
	tokens := l.Tokenize()
	var modifiers []Token
	for _, tok := range tokens {
		if tok.Type == TokenStatic || tok.Type == TokenField || tok.Type == TokenMethod {
			modifiers = append(modifiers, tok)
		}
	}
	require.Len(t, modifiers, 3)
	assert.Equal(t, TokenStatic, modifiers[0].Type)
	assert.Equal(t, TokenField, modifiers[1].Type)
	assert.Equal(t, TokenMethod, modifiers[2].Type)
}

func TestPos_String(t *testing.T) {
	t.Parallel()
	p := Pos{Line: 5, Column: 10}
	assert.Equal(t, "5:10", p.String())
}

func TestToken_String(t *testing.T) {
	t.Parallel()
	tok := Token{Type: TokenClass, Literal: "class", Pos: Pos{1, 1}}
	assert.Contains(t, tok.String(), "class")
	assert.Contains(t, tok.String(), "1:1")
}
