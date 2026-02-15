package parser

import (
	"fmt"
	"strings"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/bobcob7/go-uml/internal/lexer"
)

func (p *Parser) parseParticipant(kind ast.ParticipantKind) *ast.Participant {
	tok := p.advance() // consume keyword
	name := p.readParticipantName()
	alias := ""
	if p.current().Type == lexer.TokenAs {
		p.advance()
		if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
			alias = stripQuotes(p.current().Literal)
			p.advance()
		}
	}
	p.skipToNextLine()
	return &ast.Participant{Pos: tok.Pos, Name: name, Alias: alias, Kind: kind}
}

func (p *Parser) readParticipantName() string {
	if p.current().Type == lexer.TokenString {
		name := stripQuotes(p.current().Literal)
		p.advance()
		return name
	}
	if p.current().Type == lexer.TokenIdent {
		name := p.current().Literal
		p.advance()
		return name
	}
	return ""
}

func (p *Parser) parseActivate(deactivate bool) *ast.Activate {
	tok := p.advance() // consume 'activate' or 'deactivate'
	target := ""
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		target = stripQuotes(p.current().Literal)
		p.advance()
	}
	p.skipToNextLine()
	return &ast.Activate{Pos: tok.Pos, Target: target, Deactivate: deactivate}
}

func (p *Parser) parseReturn() *ast.Return {
	tok := p.advance() // consume 'return'
	label := p.readRestOfLine()
	return &ast.Return{Pos: tok.Pos, Label: label}
}

func (p *Parser) parseFragment(kind ast.FragmentKind) *ast.Fragment {
	tok := p.advance() // consume keyword (alt, loop, par, etc.)
	condition := p.readRestOfLine()
	frag := &ast.Fragment{
		Pos:       tok.Pos,
		Kind:      kind,
		Condition: strings.TrimSpace(condition),
	}
	p.skipNewlines()
	for p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenElse &&
		p.current().Type != lexer.TokenEndUML && p.current().Type != lexer.TokenEOF {
		p.skipNewlines()
		if p.current().Type == lexer.TokenEnd || p.current().Type == lexer.TokenElse ||
			p.current().Type == lexer.TokenEndUML || p.current().Type == lexer.TokenEOF {
			break
		}
		stmt := p.parseStatementInContext(true)
		if stmt != nil {
			frag.Statements = append(frag.Statements, stmt)
		}
	}
	for p.current().Type == lexer.TokenElse {
		ep := p.parseElsePart()
		frag.ElseParts = append(frag.ElseParts, ep)
	}
	if p.current().Type == lexer.TokenEnd {
		p.advance()
		p.skipToNextLine()
	} else {
		p.addError(p.current().Pos, fmt.Sprintf("expected 'end' to close %s fragment", fragmentName(kind)))
	}
	return frag
}

func (p *Parser) parseElsePart() ast.ElsePart {
	tok := p.advance() // consume 'else'
	condition := p.readRestOfLine()
	ep := ast.ElsePart{
		Pos:       tok.Pos,
		Condition: strings.TrimSpace(condition),
	}
	p.skipNewlines()
	for p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenElse &&
		p.current().Type != lexer.TokenEndUML && p.current().Type != lexer.TokenEOF {
		p.skipNewlines()
		if p.current().Type == lexer.TokenEnd || p.current().Type == lexer.TokenElse ||
			p.current().Type == lexer.TokenEndUML || p.current().Type == lexer.TokenEOF {
			break
		}
		stmt := p.parseStatementInContext(true)
		if stmt != nil {
			ep.Statements = append(ep.Statements, stmt)
		}
	}
	return ep
}

func (p *Parser) parseAutonumber() *ast.Autonumber {
	tok := p.advance() // consume 'autonumber'
	start := ""
	if p.current().Type == lexer.TokenNumber {
		start = p.current().Literal
		p.advance()
	}
	p.skipToNextLine()
	return &ast.Autonumber{Pos: tok.Pos, Start: start}
}

func (p *Parser) parseDivider() *ast.Divider {
	tok := p.advance() // consume first '='
	if p.current().Type == lexer.TokenEquals {
		p.advance()
	}
	var parts []string
	for p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		if p.current().Type == lexer.TokenEquals {
			p.advance()
			if p.current().Type == lexer.TokenEquals {
				p.advance()
				break
			}
			parts = append(parts, "=")
			continue
		}
		parts = append(parts, p.current().Literal)
		p.advance()
	}
	text := strings.TrimSpace(strings.Join(parts, " "))
	return &ast.Divider{Pos: tok.Pos, Text: text}
}

func (p *Parser) parseDelay() *ast.Delay {
	tok := p.advance() // consume '...' arrow
	text := ""
	if p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
		var parts []string
		for p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenEOF {
			if p.current().Type == lexer.TokenArrow && isDelayArrow(p.current().Literal) {
				p.advance()
				break
			}
			parts = append(parts, p.current().Literal)
			p.advance()
		}
		text = strings.TrimSpace(strings.Join(parts, " "))
	}
	return &ast.Delay{Pos: tok.Pos, Text: text}
}

func (p *Parser) parseSequenceIdentStatement() ast.Statement {
	tok := p.current()
	name := tok.Literal
	p.advance()
	if p.current().Type == lexer.TokenArrow {
		return p.parseMessage(tok.Pos, name)
	}
	p.addError(tok.Pos, fmt.Sprintf("unexpected identifier %q", name))
	p.skipToNextLine()
	return nil
}

func (p *Parser) parseMessage(pos lexer.Pos, from string) *ast.Message {
	arrowTok := p.advance() // consume arrow
	arrow := arrowTok.Literal
	dashed := isDashedArrow(arrow)
	to := ""
	if p.current().Type == lexer.TokenIdent || p.current().Type == lexer.TokenString {
		to = stripQuotes(p.current().Literal)
		p.advance()
	}
	// Handle activation shorthand: ++ or --
	activate := ""
	if p.current().Type == lexer.TokenPlus {
		p.advance()
		if p.current().Type == lexer.TokenPlus {
			p.advance()
			activate = "++"
		}
	} else if p.current().Type == lexer.TokenArrow && p.current().Literal == "--" {
		p.advance()
		activate = "--"
	}
	label := ""
	if p.current().Type == lexer.TokenColon {
		p.advance()
		label = strings.TrimSpace(p.readRestOfLine())
	} else {
		p.skipToNextLine()
	}
	msg := &ast.Message{
		Pos:    pos,
		From:   from,
		To:     to,
		Label:  label,
		Arrow:  arrow,
		Dashed: dashed,
	}
	_ = activate // activation shorthand tracked but not yet wired to AST
	return msg
}

func isDashedArrow(arrow string) bool {
	shaft := strings.TrimLeft(arrow, "<|")
	shaft = strings.TrimRight(shaft, ">|*o")
	return strings.Contains(shaft, "..") || strings.Contains(shaft, "--")
}

func isDelayArrow(literal string) bool {
	for _, ch := range literal {
		if ch != '.' {
			return false
		}
	}
	return len(literal) >= 3
}

func fragmentName(kind ast.FragmentKind) string {
	switch kind {
	case ast.FragmentAlt:
		return "alt"
	case ast.FragmentLoop:
		return "loop"
	case ast.FragmentPar:
		return "par"
	case ast.FragmentBreak:
		return "break"
	case ast.FragmentRef:
		return "ref"
	case ast.FragmentGroup:
		return "group"
	default:
		return "unknown"
	}
}

func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
