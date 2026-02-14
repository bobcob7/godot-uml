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
	case lexer.TokenNote:
		return p.parseNote()
	case lexer.TokenIdent:
		return p.parseIdentStatement()
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

func (p *Parser) parseAbstract() ast.Statement {
	tok := p.advance() // consume 'abstract'
	if p.current().Type == lexer.TokenClass {
		return p.parseClassDef(true)
	}
	// "abstract" alone is treated as an abstract class with the next identifier as name.
	cd := &ast.ClassDef{Pos: tok.Pos, Abstract: true}
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		cd.Name = p.readClassName()
	}
	cd.Stereotype = p.tryStereotype()
	if p.current().Type == lexer.TokenLBrace {
		cd.Members = p.parseClassBody()
	}
	return cd
}

func (p *Parser) parseClassDef(abstract bool) *ast.ClassDef {
	tok := p.advance() // consume 'class'
	if !abstract {
		tok.Pos = lexer.Pos{Line: tok.Pos.Line, Column: tok.Pos.Column}
	}
	cd := &ast.ClassDef{Pos: tok.Pos, Abstract: abstract}
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		cd.Name = p.readClassName()
	} else {
		p.addError(p.current().Pos, "expected class name")
		p.skipToNextLine()
		return cd
	}
	cd.Stereotype = p.tryStereotype()
	if p.current().Type == lexer.TokenAs {
		p.advance()
		if p.current().Type == lexer.TokenIdent {
			cd.Alias = p.current().Literal
			p.advance()
		}
	}
	if p.current().Type == lexer.TokenExtends || p.current().Type == lexer.TokenImplements {
		p.skipToNextLine()
		return cd
	}
	if p.current().Type == lexer.TokenLBrace {
		cd.Members = p.parseClassBody()
	}
	return cd
}

func (p *Parser) parseInterfaceDef() *ast.InterfaceDef {
	tok := p.advance() // consume 'interface'
	idef := &ast.InterfaceDef{Pos: tok.Pos}
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		idef.Name = p.readClassName()
	} else {
		p.addError(p.current().Pos, "expected interface name")
		p.skipToNextLine()
		return idef
	}
	idef.Stereotype = p.tryStereotype()
	if p.current().Type == lexer.TokenAs {
		p.advance()
		if p.current().Type == lexer.TokenIdent {
			idef.Alias = p.current().Literal
			p.advance()
		}
	}
	if p.current().Type == lexer.TokenLBrace {
		idef.Members = p.parseClassBody()
	}
	return idef
}

func (p *Parser) parseEnumDef() *ast.EnumDef {
	tok := p.advance() // consume 'enum'
	edef := &ast.EnumDef{Pos: tok.Pos}
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		edef.Name = p.readClassName()
	} else {
		p.addError(p.current().Pos, "expected enum name")
		p.skipToNextLine()
		return edef
	}
	edef.Stereotype = p.tryStereotype()
	if p.current().Type == lexer.TokenLBrace {
		edef.Members = p.parseClassBody()
	}
	return edef
}

// readClassName reads a class name, which may be a dotted identifier (e.g. "com.example.Foo").
func (p *Parser) readClassName() string {
	if p.current().Type == lexer.TokenString {
		tok := p.advance()
		return strings.Trim(tok.Literal, "\"")
	}
	var b strings.Builder
	b.WriteString(p.current().Literal)
	p.advance()
	for p.current().Type == lexer.TokenDot && p.peek().Type == lexer.TokenIdent {
		b.WriteRune('.')
		p.advance() // consume dot
		b.WriteString(p.current().Literal)
		p.advance() // consume ident
	}
	return b.String()
}

// tryStereotype checks for <<stereotype>> and returns the text, or "" if none.
func (p *Parser) tryStereotype() string {
	if p.current().Type != lexer.TokenLAngle {
		return ""
	}
	if p.peek().Type != lexer.TokenLAngle {
		return ""
	}
	p.advance() // first <
	p.advance() // second <
	var parts []string
	for p.current().Type != lexer.TokenRAngle && p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		parts = append(parts, p.current().Literal)
		p.advance()
	}
	if p.current().Type == lexer.TokenRAngle {
		p.advance() // first >
	}
	if p.current().Type == lexer.TokenRAngle {
		p.advance() // second >
	}
	return strings.Join(parts, " ")
}

func (p *Parser) parseClassBody() []ast.Member {
	p.advance() // consume '{'
	p.skipNewlines()
	var members []ast.Member
	for p.current().Type != lexer.TokenRBrace && p.current().Type != lexer.TokenEOF {
		p.skipNewlines()
		if p.current().Type == lexer.TokenRBrace || p.current().Type == lexer.TokenEOF {
			break
		}
		if p.current().Type == lexer.TokenLineComment || p.current().Type == lexer.TokenBlockComment {
			p.advance()
			continue
		}
		member := p.parseMember()
		if member != nil {
			members = append(members, member)
		}
	}
	if p.current().Type == lexer.TokenRBrace {
		p.advance()
	} else {
		p.addError(p.current().Pos, "expected closing }")
	}
	return members
}

func (p *Parser) parseMember() ast.Member {
	pos := p.current().Pos
	vis := p.tryVisibility()
	mod := p.tryModifier()
	if p.current().Type == lexer.TokenNewline || p.current().Type == lexer.TokenRBrace || p.current().Type == lexer.TokenEOF {
		return nil
	}
	name := ""
	if p.current().Type == lexer.TokenIdent {
		name = p.current().Literal
		p.advance()
	} else {
		// Consume the rest of the line as a member name.
		name = p.readRestOfLine()
		if name == "" {
			p.skipToNextLine()
			return nil
		}
		return &ast.Field{Pos: pos, Name: name, Visibility: vis, Modifier: mod}
	}
	// Check if it's a method (has parentheses).
	if p.current().Type == lexer.TokenLParen {
		return p.parseMethodAfterName(pos, vis, mod, name)
	}
	// It's a field. Check for type annotation.
	typeName := ""
	if p.current().Type == lexer.TokenColon {
		p.advance()
		typeName = p.readTypeUntilNewline()
	}
	p.consumeOptionalNewline()
	return &ast.Field{Pos: pos, Name: name, Type: typeName, Visibility: vis, Modifier: mod}
}

func (p *Parser) parseMethodAfterName(pos lexer.Pos, vis ast.Visibility, mod ast.Modifier, name string) *ast.Method {
	p.advance() // consume '('
	var params []string
	for p.current().Type != lexer.TokenRParen && p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		params = append(params, p.current().Literal)
		p.advance()
	}
	if p.current().Type == lexer.TokenRParen {
		p.advance()
	}
	retType := ""
	if p.current().Type == lexer.TokenColon {
		p.advance()
		retType = p.readTypeUntilNewline()
	}
	p.consumeOptionalNewline()
	return &ast.Method{
		Pos:        pos,
		Name:       name,
		Params:     strings.Join(params, " "),
		ReturnType: retType,
		Visibility: vis,
		Modifier:   mod,
	}
}

func (p *Parser) tryVisibility() ast.Visibility {
	switch p.current().Type {
	case lexer.TokenPlus:
		p.advance()
		return ast.VisibilityPublic
	case lexer.TokenMinus:
		p.advance()
		return ast.VisibilityPrivate
	case lexer.TokenHash:
		p.advance()
		return ast.VisibilityProtected
	case lexer.TokenTilde:
		p.advance()
		return ast.VisibilityPackage
	default:
		return ast.VisibilityNone
	}
}

func (p *Parser) tryModifier() ast.Modifier {
	switch p.current().Type {
	case lexer.TokenStatic:
		p.advance()
		return ast.ModifierStatic
	case lexer.TokenField:
		p.advance()
		return ast.ModifierField
	case lexer.TokenMethod:
		p.advance()
		return ast.ModifierMethod
	default:
		return ast.ModifierNone
	}
}

func (p *Parser) readTypeUntilNewline() string {
	var parts []string
	for p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF &&
		p.current().Type != lexer.TokenRBrace {
		parts = append(parts, p.current().Literal)
		p.advance()
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (p *Parser) consumeOptionalNewline() {
	if p.current().Type == lexer.TokenNewline {
		p.advance()
	}
}

// parseIdentStatement handles lines starting with an identifier, which could be
// a relationship (e.g., "Foo --> Bar") or other identifier-based statement.
func (p *Parser) parseIdentStatement() ast.Statement {
	pos := p.current().Pos
	leftName := p.readClassName()
	// Check for cardinality before arrow: Foo "1" --> "*" Bar
	leftCard := ""
	if p.current().Type == lexer.TokenString {
		leftCard = strings.Trim(p.current().Literal, "\"")
		p.advance()
	}
	if p.current().Type == lexer.TokenArrow {
		return p.parseRelationship(pos, leftName, leftCard)
	}
	// Not a relationship, skip rest of line.
	p.skipToNextLine()
	return &ast.Comment{Pos: pos, Text: leftName}
}

func (p *Parser) parseRelationship(pos lexer.Pos, leftName, leftCard string) *ast.Relationship {
	arrowTok := p.advance() // consume arrow
	relType, dir := classifyArrow(arrowTok.Literal)
	rightCard := ""
	if p.current().Type == lexer.TokenString {
		rightCard = strings.Trim(p.current().Literal, "\"")
		p.advance()
	}
	rightName := ""
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		rightName = p.readClassName()
	}
	label := ""
	if p.current().Type == lexer.TokenColon {
		p.advance()
		label = strings.TrimSpace(p.readRestOfLine())
	}
	return &ast.Relationship{
		Pos:       pos,
		Left:      leftName,
		Right:     rightName,
		Type:      relType,
		Direction: dir,
		Label:     label,
		LeftCard:  leftCard,
		RightCard: rightCard,
		Arrow:     arrowTok.Literal,
	}
}

// classifyArrow determines the relationship type and direction from an arrow literal.
func classifyArrow(arrow string) (ast.RelationshipType, ast.ArrowDirection) {
	dir := arrowDirection(arrow)
	switch {
	case strings.Contains(arrow, "|>") || strings.Contains(arrow, "<|"):
		if strings.Contains(arrow, "..") {
			return ast.RelRealization, dir
		}
		return ast.RelInheritance, dir
	case strings.HasSuffix(arrow, "*") || strings.HasPrefix(arrow, "*"):
		return ast.RelComposition, dir
	case (strings.HasSuffix(arrow, "o") && len(arrow) > 1) || strings.HasPrefix(arrow, "o"):
		return ast.RelAggregation, dir
	case strings.Contains(arrow, ".."):
		if strings.Contains(arrow, ">") || strings.Contains(arrow, "<") {
			return ast.RelDependency, dir
		}
		return ast.RelDependency, dir
	case strings.Contains(arrow, ">") || strings.Contains(arrow, "<"):
		return ast.RelAssociation, dir
	default:
		return ast.RelAssociation, dir
	}
}

func arrowDirection(arrow string) ast.ArrowDirection {
	hasLeft := strings.HasPrefix(arrow, "<")
	hasRight := strings.HasSuffix(arrow, ">") || strings.HasSuffix(arrow, "|>")
	switch {
	case hasLeft && hasRight:
		return ast.ArrowBoth
	case hasLeft:
		return ast.ArrowLeft
	case hasRight:
		return ast.ArrowRight
	default:
		return ast.ArrowNone
	}
}

func (p *Parser) parsePackage() ast.Statement {
	tok := p.advance() // consume 'package' or 'namespace'
	isNamespace := tok.Type == lexer.TokenNamespace
	pkg := &ast.Package{Pos: tok.Pos, IsNamespace: isNamespace}
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		pkg.Name = p.readClassName()
	}
	pkg.Alias = ""
	if p.current().Type == lexer.TokenAs {
		p.advance()
		if p.current().Type == lexer.TokenIdent {
			pkg.Alias = p.current().Literal
			p.advance()
		}
	}
	if p.current().Type == lexer.TokenLBrace {
		p.advance()
		p.skipNewlines()
		for p.current().Type != lexer.TokenRBrace && p.current().Type != lexer.TokenEOF {
			p.skipNewlines()
			if p.current().Type == lexer.TokenRBrace || p.current().Type == lexer.TokenEOF {
				break
			}
			stmt := p.parseStatement()
			if stmt != nil {
				pkg.Statements = append(pkg.Statements, stmt)
			}
		}
		if p.current().Type == lexer.TokenRBrace {
			p.advance()
		} else {
			p.addError(p.current().Pos, "expected closing } for package")
		}
	}
	return pkg
}

func (p *Parser) parseNote() ast.Statement {
	tok := p.advance() // consume 'note'
	n := &ast.Note{Pos: tok.Pos}
	switch p.current().Type {
	case lexer.TokenLeft:
		n.Placement = ast.NoteLeft
		p.advance()
	case lexer.TokenRight:
		n.Placement = ast.NoteRight
		p.advance()
	case lexer.TokenOver:
		n.Placement = ast.NoteOver
		p.advance()
	}
	if p.current().Type == lexer.TokenOf {
		p.advance()
	}
	if p.current().Type == lexer.TokenIdent {
		n.Target = p.current().Literal
		p.advance()
	}
	if p.current().Type == lexer.TokenColon {
		p.advance()
		n.Text = strings.TrimSpace(p.readRestOfLine())
	} else {
		n.Text = strings.TrimSpace(p.readRestOfLine())
	}
	return n
}
