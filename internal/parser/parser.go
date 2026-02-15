// Package parser implements a recursive descent parser that converts tokens into an AST.
package parser

import (
	"fmt"
	"strings"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/bobcob7/go-uml/internal/lexer"
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
	tokens  []lexer.Token
	pos     int
	errors  []*Error
	seqMode bool // true after a sequence-specific keyword is seen
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

func (p *Parser) peek() lexer.Token {
	next := p.pos + 1
	if next >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[next]
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
	return p.parseStatementInContext(false)
}

func (p *Parser) parseStatementInContext(inFragment bool) ast.Statement {
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
	case lexer.TokenClass:
		return p.parseClassDef(false)
	case lexer.TokenAbstract:
		return p.parseAbstract()
	case lexer.TokenInterface:
		return p.parseInterfaceDef()
	case lexer.TokenEnum:
		return p.parseEnumDef()
	case lexer.TokenPackage, lexer.TokenNamespace:
		return p.parsePackage()
	case lexer.TokenParticipant:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantDefault)
	case lexer.TokenActor:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantActor)
	case lexer.TokenBoundary:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantBoundary)
	case lexer.TokenControl:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantControl)
	case lexer.TokenEntity:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantEntity)
	case lexer.TokenDatabase:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantDatabase)
	case lexer.TokenCollections:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantCollections)
	case lexer.TokenQueue:
		p.seqMode = true
		return p.parseParticipant(ast.ParticipantQueue)
	case lexer.TokenActivate:
		p.seqMode = true
		return p.parseActivate(false)
	case lexer.TokenDeactivate:
		p.seqMode = true
		return p.parseActivate(true)
	case lexer.TokenReturn:
		p.seqMode = true
		return p.parseReturn()
	case lexer.TokenAlt:
		p.seqMode = true
		return p.parseFragment(ast.FragmentAlt)
	case lexer.TokenLoop:
		p.seqMode = true
		return p.parseFragment(ast.FragmentLoop)
	case lexer.TokenPar:
		p.seqMode = true
		return p.parseFragment(ast.FragmentPar)
	case lexer.TokenBreak:
		p.seqMode = true
		return p.parseFragment(ast.FragmentBreak)
	case lexer.TokenRef:
		p.seqMode = true
		return p.parseFragment(ast.FragmentRef)
	case lexer.TokenGroup:
		p.seqMode = true
		return p.parseFragment(ast.FragmentGroup)
	case lexer.TokenAutonumber:
		p.seqMode = true
		return p.parseAutonumber()
	case lexer.TokenNote:
		return p.parseNote()
	case lexer.TokenEquals:
		return p.parseDivider()
	case lexer.TokenArrow:
		if isDelayArrow(tok.Literal) {
			return p.parseDelay()
		}
		p.addError(tok.Pos, fmt.Sprintf("unexpected arrow %q", tok.Literal))
		p.skipToNextLine()
		return nil
	case lexer.TokenIdent:
		return p.parseIdentStatement()
	case lexer.TokenError:
		p.addError(tok.Pos, fmt.Sprintf("unexpected token: %s", tok.Literal))
		p.skipToNextLine()
		return nil
	default:
		if inFragment && (tok.Type == lexer.TokenElse || tok.Type == lexer.TokenEnd) {
			return nil
		}
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
	tok := p.advance()
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "title " + text}
}

func (p *Parser) parseHeader() ast.Statement {
	tok := p.advance()
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "header " + text}
}

func (p *Parser) parseFooter() ast.Statement {
	tok := p.advance()
	text := p.readRestOfLine()
	return &ast.Comment{Pos: tok.Pos, Text: "footer " + text}
}

func (p *Parser) parseSkinparam() *ast.Skinparam {
	tok := p.advance()
	name := ""
	if p.current().Type == lexer.TokenIdent {
		name = p.current().Literal
		p.advance()
	}
	value := p.readRestOfLine()
	return &ast.Skinparam{Pos: tok.Pos, Name: name, Value: value}
}

func (p *Parser) parseHideShow(isHide bool) *ast.HideShow {
	tok := p.advance()
	target := p.readRestOfLine()
	return &ast.HideShow{Pos: tok.Pos, IsHide: isHide, Target: target}
}

func (p *Parser) parseNote() *ast.Note {
	tok := p.advance() // consume 'note'
	placement := ast.NoteOver
	target := ""
	switch p.current().Type {
	case lexer.TokenLeft:
		placement = ast.NoteLeft
		p.advance()
		if p.current().Type == lexer.TokenOf {
			p.advance()
		}
		target = p.readNoteTarget()
	case lexer.TokenRight:
		placement = ast.NoteRight
		p.advance()
		if p.current().Type == lexer.TokenOf {
			p.advance()
		}
		target = p.readNoteTarget()
	case lexer.TokenOver:
		p.advance()
		target = p.readNoteTarget()
	}
	text := ""
	if p.current().Type == lexer.TokenColon {
		p.advance()
		text = strings.TrimSpace(p.readRestOfLine())
	} else {
		text = p.readMultiLineNote()
	}
	return &ast.Note{Pos: tok.Pos, Placement: placement, Target: target, Text: text}
}

func (p *Parser) readNoteTarget() string {
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		name := stripQuotes(p.current().Literal)
		p.advance()
		for p.current().Type == lexer.TokenComma {
			name += ","
			p.advance()
			if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
				name += stripQuotes(p.current().Literal)
				p.advance()
			}
		}
		return name
	}
	return ""
}

func (p *Parser) readMultiLineNote() string {
	var lines []string
	for p.current().Type != lexer.TokenEOF && p.current().Type != lexer.TokenEndUML {
		if p.current().Type == lexer.TokenNewline {
			p.advance()
			continue
		}
		if p.current().Type == lexer.TokenEnd {
			p.advance()
			if p.current().Type == lexer.TokenNote {
				p.advance()
			}
			p.skipToNextLine()
			break
		}
		lines = append(lines, p.readRestOfLine())
	}
	return strings.Join(lines, "\n")
}
