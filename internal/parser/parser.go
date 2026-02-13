// Package parser implements a recursive descent parser that converts tokens into an AST.
package parser

import (
	"fmt"
	"strings"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
)

// Error represents a parse error with source position.
type Error struct {
	Pos     lexer.Pos
	Message string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Pos, e.Message)
}

// Parser is a recursive descent parser for PlantUML diagrams.
type Parser struct {
	tokens []lexer.Token
	pos    int
	errors []*Error
}

// New creates a new Parser for the given token slice.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// Parse parses a complete PlantUML diagram and returns the AST root plus any errors.
// The parser uses error recovery to continue after errors and report multiple issues.
func Parse(input string) (*ast.Diagram, []*Error) {
	l := lexer.New(input)
	tokens := l.Tokenize()
	p := New(tokens)
	diagram := p.parseDiagram()
	return diagram, p.errors
}

// Errors returns all parse errors collected during parsing.
func (p *Parser) Errors() []*Error {
	return p.errors
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *Parser) addError(pos lexer.Pos, msg string) {
	p.errors = append(p.errors, &Error{Pos: pos, Message: msg})
}

func (p *Parser) skipNewlines() {
	for p.current().Type == lexer.TokenNewline {
		p.advance()
	}
}

// skipToNextLine advances past all tokens until the next newline or EOF, used for error recovery.
func (p *Parser) skipToNextLine() {
	for p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		p.advance()
	}
	if p.current().Type == lexer.TokenNewline {
		p.advance()
	}
}

// readRestOfLine reads remaining tokens on the current line as space-joined text.
func (p *Parser) readRestOfLine() string {
	var parts []string
	for p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		parts = append(parts, p.current().Literal)
		p.advance()
	}
	return strings.Join(parts, " ")
}

func (p *Parser) parseDiagram() *ast.Diagram {
	p.skipNewlines()
	diagram := &ast.Diagram{}
	tok := p.current()
	switch tok.Type {
	case lexer.TokenStartUML:
		diagram.Pos = tok.Pos
		p.advance()
		// Optional diagram name on same line.
		if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
			diagram.Name = p.current().Literal
			p.advance()
		}
		p.skipNewlines()
	case lexer.TokenEOF:
		return diagram
	default:
		p.addError(tok.Pos, fmt.Sprintf("expected @startuml, got %s", tok.Type))
		p.skipToNextLine()
	}
	// Parse body statements.
	for p.current().Type != lexer.TokenEndUML && p.current().Type != lexer.TokenEOF {
		p.skipNewlines()
		if p.current().Type == lexer.TokenEndUML || p.current().Type == lexer.TokenEOF {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			diagram.Statements = append(diagram.Statements, stmt)
		}
	}
	if p.current().Type == lexer.TokenEndUML {
		p.advance()
	} else if p.current().Type == lexer.TokenEOF {
		p.addError(p.current().Pos, "expected @enduml before end of input")
	}
	return diagram
}

func (p *Parser) parseStatement() ast.Statement {
	tok := p.current()
	switch tok.Type {
	case lexer.TokenLineComment:
		return p.parseComment()
	case lexer.TokenBlockComment:
		return p.parseBlockComment()
	case lexer.TokenTitle:
		return p.parseTitle()
	case lexer.TokenHeader:
		return p.parseHeader()
	case lexer.TokenFooter:
		return p.parseFooter()
	case lexer.TokenSkinparam:
		return p.parseSkinparam()
	case lexer.TokenHide:
		return p.parseHideShow(true)
	case lexer.TokenShow:
		return p.parseHideShow(false)
	case lexer.TokenError:
		p.addError(tok.Pos, fmt.Sprintf("unexpected token: %s", tok.Literal))
		p.skipToNextLine()
		return nil
	default:
		p.addError(tok.Pos, fmt.Sprintf("unexpected %s %q", tok.Type, tok.Literal))
		p.skipToNextLine()
		return nil
	}
}

func (p *Parser) parseComment() ast.Statement {
	tok := p.advance()
	return &ast.Comment{Pos: tok.Pos, Text: tok.Literal}
}

func (p *Parser) parseBlockComment() ast.Statement {
	tok := p.advance()
	return &ast.Comment{Pos: tok.Pos, Text: tok.Literal}
}

func (p *Parser) parseTitle() ast.Statement {
	tok := p.advance() // consume 'title'
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "title " + text}
}

func (p *Parser) parseHeader() ast.Statement {
	tok := p.advance() // consume 'header'
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "header " + text}
}

func (p *Parser) parseFooter() ast.Statement {
	tok := p.advance() // consume 'footer'
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "footer " + text}
}

func (p *Parser) parseSkinparam() *ast.Skinparam {
	tok := p.advance() // consume 'skinparam'
	name := ""
	if p.current().Type == lexer.TokenIdent {
		name = p.current().Literal
		p.advance()
	}
	value := p.readRestOfLine()
	return &ast.Skinparam{Pos: tok.Pos, Name: name, Value: value}
}

func (p *Parser) parseHideShow(isHide bool) *ast.HideShow {
	tok := p.advance() // consume 'hide' or 'show'
	target := p.readRestOfLine()
	return &ast.HideShow{Pos: tok.Pos, IsHide: isHide, Target: target}
}
