// Package lexer tokenizes PlantUML text into a stream of typed tokens.
package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// keywords maps PlantUML keyword strings to their token types.
var keywords = map[string]TokenType{
	"class":       TokenClass,
	"interface":   TokenInterface,
	"enum":        TokenEnum,
	"abstract":    TokenAbstract,
	"extends":     TokenExtends,
	"implements":  TokenImplements,
	"package":     TokenPackage,
	"namespace":   TokenNamespace,
	"as":          TokenAs,
	"participant": TokenParticipant,
	"actor":       TokenActor,
	"boundary":    TokenBoundary,
	"control":     TokenControl,
	"entity":      TokenEntity,
	"database":    TokenDatabase,
	"collections": TokenCollections,
	"queue":       TokenQueue,
	"activate":    TokenActivate,
	"deactivate":  TokenDeactivate,
	"return":      TokenReturn,
	"alt":         TokenAlt,
	"else":        TokenElse,
	"end":         TokenEnd,
	"loop":        TokenLoop,
	"group":       TokenGroup,
	"note":        TokenNote,
	"of":          TokenOf,
	"over":        TokenOver,
	"left":        TokenLeft,
	"right":       TokenRight,
	"skinparam":   TokenSkinparam,
	"hide":        TokenHide,
	"show":        TokenShow,
	"title":       TokenTitle,
	"header":      TokenHeader,
	"footer":      TokenFooter,
}

// Lexer tokenizes PlantUML source text.
type Lexer struct {
	input string
	pos   int  // current byte position in input
	line  int  // current line (1-based)
	col   int  // current column (1-based)
	ch    rune // current character
	eof   bool // true when past end of input
}

// New creates a new Lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   1,
	}
	l.readChar()
	return l
}

// Tokenize consumes the entire input and returns all tokens including the final EOF.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	if l.eof {
		return Token{Type: TokenEOF, Pos: l.currentPos()}
	}
	pos := l.currentPos()
	switch {
	case l.ch == '@':
		return l.readDirective(pos)
	case l.ch == '\n':
		l.readChar()
		return Token{Type: TokenNewline, Literal: "\n", Pos: pos}
	case l.ch == '{':
		return l.readBraceOrModifier(pos)
	case l.ch == '}':
		l.readChar()
		return Token{Type: TokenRBrace, Literal: "}", Pos: pos}
	case l.ch == '(':
		l.readChar()
		return Token{Type: TokenLParen, Literal: "(", Pos: pos}
	case l.ch == ')':
		l.readChar()
		return Token{Type: TokenRParen, Literal: ")", Pos: pos}
	case l.ch == '[':
		l.readChar()
		return Token{Type: TokenLBracket, Literal: "[", Pos: pos}
	case l.ch == ']':
		l.readChar()
		return Token{Type: TokenRBracket, Literal: "]", Pos: pos}
	case l.ch == ':':
		l.readChar()
		return Token{Type: TokenColon, Literal: ":", Pos: pos}
	case l.ch == ',':
		l.readChar()
		return Token{Type: TokenComma, Literal: ",", Pos: pos}
	case l.ch == '.':
		if l.peekChar() == '.' {
			return l.readArrowOrDots(pos)
		}
		l.readChar()
		return Token{Type: TokenDot, Literal: ".", Pos: pos}
	case l.ch == '|':
		l.readChar()
		return Token{Type: TokenPipe, Literal: "|", Pos: pos}
	case l.ch == '=':
		l.readChar()
		return Token{Type: TokenEquals, Literal: "=", Pos: pos}
	case l.ch == ';':
		l.readChar()
		return Token{Type: TokenSemicolon, Literal: ";", Pos: pos}
	case l.ch == '+':
		l.readChar()
		return Token{Type: TokenPlus, Literal: "+", Pos: pos}
	case l.ch == '~':
		l.readChar()
		return Token{Type: TokenTilde, Literal: "~", Pos: pos}
	case l.ch == '#':
		l.readChar()
		return Token{Type: TokenHash, Literal: "#", Pos: pos}
	case l.ch == '\'':
		return l.readComment(pos)
	case l.ch == '/':
		if l.peekChar() == '\'' {
			return l.readBlockComment(pos)
		}
		return l.readError(pos)
	case l.ch == '"':
		return l.readString('"', pos)
	case l.ch == '-':
		return l.readDashStart(pos)
	case l.ch == '<':
		return l.readLeftAngleStart(pos)
	case l.ch == '*':
		return l.readStarStart(pos)
	case l.ch == 'o':
		if l.isArrowContinuation(l.peekChar()) {
			return l.readArrowFrom(pos)
		}
		return l.readIdentOrKeyword(pos)
	case l.ch == '>':
		l.readChar()
		return Token{Type: TokenRAngle, Literal: ">", Pos: pos}
	case unicode.IsLetter(l.ch) || l.ch == '_':
		return l.readIdentOrKeyword(pos)
	case unicode.IsDigit(l.ch):
		return l.readNumber(pos)
	default:
		return l.readError(pos)
	}
}

func (l *Lexer) currentPos() Pos {
	return Pos{Line: l.line, Column: l.col}
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
		l.eof = true
		return
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	if l.ch == '\n' {
		l.line++
		l.col = 1
	} else if !l.eof {
		l.col++
	}
	if l.pos == 0 {
		l.col = 1
	}
	l.ch = r
	l.pos += size
}

func (l *Lexer) peekChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *Lexer) skipWhitespace() {
	for !l.eof && l.ch != '\n' && unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readDirective(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // @
	l.readChar()
	for !l.eof && (unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch)) {
		b.WriteRune(l.ch)
		l.readChar()
	}
	lit := b.String()
	switch strings.ToLower(lit) {
	case "@startuml":
		return Token{Type: TokenStartUML, Literal: lit, Pos: pos}
	case "@enduml":
		return Token{Type: TokenEndUML, Literal: lit, Pos: pos}
	default:
		return Token{Type: TokenError, Literal: lit, Pos: pos}
	}
}

func (l *Lexer) readBraceOrModifier(pos Pos) Token {
	// Check for {static}, {field}, {method} modifiers.
	rest := l.input[l.pos:] // text after '{'
	for _, kw := range []struct {
		text string
		tt   TokenType
	}{
		{"static}", TokenStatic},
		{"field}", TokenField},
		{"method}", TokenMethod},
	} {
		if strings.HasPrefix(rest, kw.text) {
			lit := "{" + kw.text
			for range len(lit) - 1 { // -1 because we haven't consumed '{' yet... consume all chars
				l.readChar()
			}
			l.readChar() // consume last char
			return Token{Type: kw.tt, Literal: lit, Pos: pos}
		}
	}
	l.readChar()
	return Token{Type: TokenLBrace, Literal: "{", Pos: pos}
}

func (l *Lexer) readComment(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // '
	l.readChar()
	for !l.eof && l.ch != '\n' {
		b.WriteRune(l.ch)
		l.readChar()
	}
	return Token{Type: TokenLineComment, Literal: b.String(), Pos: pos}
}

func (l *Lexer) readBlockComment(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // /
	l.readChar()      // now at '
	b.WriteRune(l.ch)
	l.readChar()
	for !l.eof {
		if l.ch == '\'' && l.peekChar() == '/' {
			b.WriteRune(l.ch)
			l.readChar()
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenBlockComment, Literal: b.String(), Pos: pos}
		}
		b.WriteRune(l.ch)
		l.readChar()
	}
	// Unterminated block comment.
	return Token{Type: TokenError, Literal: b.String(), Pos: pos}
}

func (l *Lexer) readString(quote rune, pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // opening quote
	l.readChar()
	for !l.eof && l.ch != quote && l.ch != '\n' {
		if l.ch == '\\' {
			b.WriteRune(l.ch)
			l.readChar()
			if !l.eof {
				b.WriteRune(l.ch)
				l.readChar()
			}
			continue
		}
		b.WriteRune(l.ch)
		l.readChar()
	}
	if l.ch == quote {
		b.WriteRune(l.ch)
		l.readChar()
		return Token{Type: TokenString, Literal: b.String(), Pos: pos}
	}
	return Token{Type: TokenError, Literal: b.String(), Pos: pos}
}

func (l *Lexer) readIdentOrKeyword(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch)
	l.readChar()
	for !l.eof && (unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_') {
		b.WriteRune(l.ch)
		l.readChar()
	}
	lit := b.String()
	if tt, ok := keywords[lit]; ok {
		return Token{Type: tt, Literal: lit, Pos: pos}
	}
	return Token{Type: TokenIdent, Literal: lit, Pos: pos}
}

func (l *Lexer) readNumber(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch)
	l.readChar()
	for !l.eof && (unicode.IsDigit(l.ch) || l.ch == '.') {
		b.WriteRune(l.ch)
		l.readChar()
	}
	return Token{Type: TokenNumber, Literal: b.String(), Pos: pos}
}

func (l *Lexer) readError(pos Pos) Token {
	lit := string(l.ch)
	l.readChar()
	return Token{Type: TokenError, Literal: lit, Pos: pos}
}

// Arrow reading helpers.

// isArrowContinuation returns true if r can follow 'o', '*', or '<' in an arrow.
func (l *Lexer) isArrowContinuation(r rune) bool {
	return r == '-' || r == '.'
}

// readArrowFrom reads an arrow starting with a collected prefix (e.g. "<|", "*", "o").
func (l *Lexer) readArrowFrom(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // first char (e.g. 'o', '*', '<')
	l.readChar()
	return l.continueArrow(&b, pos)
}

// readDashStart handles tokens starting with '-'.
func (l *Lexer) readDashStart(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // -
	l.readChar()
	// Consume consecutive dashes/dots to form arrow shaft.
	for !l.eof && (l.ch == '-' || l.ch == '.') {
		b.WriteRune(l.ch)
		l.readChar()
	}
	// Check for arrowhead at end: >, |>, >>, etc.
	if l.eof {
		return l.finishArrowOrMinus(&b, pos)
	}
	switch l.ch {
	case '>':
		b.WriteRune(l.ch)
		l.readChar()
		return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
	case '|':
		b.WriteRune(l.ch)
		l.readChar()
		if !l.eof && l.ch == '>' {
			b.WriteRune(l.ch)
			l.readChar()
		}
		return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
	case '*':
		b.WriteRune(l.ch)
		l.readChar()
		return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
	case 'o':
		b.WriteRune(l.ch)
		l.readChar()
		return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
	default:
		return l.finishArrowOrMinus(&b, pos)
	}
}

func (l *Lexer) finishArrowOrMinus(b *strings.Builder, pos Pos) Token {
	lit := b.String()
	if len(lit) >= 2 {
		// Bare shaft like "--" is still an arrow (e.g. association).
		return Token{Type: TokenArrow, Literal: lit, Pos: pos}
	}
	return Token{Type: TokenMinus, Literal: lit, Pos: pos}
}

// readLeftAngleStart handles tokens starting with '<'.
func (l *Lexer) readLeftAngleStart(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch) // <
	l.readChar()
	switch l.ch {
	case '|':
		b.WriteRune(l.ch)
		l.readChar()
		if l.ch == '-' || l.ch == '.' {
			return l.continueArrow(&b, pos)
		}
		return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
	case '-', '.':
		return l.continueArrow(&b, pos)
	default:
		return Token{Type: TokenLAngle, Literal: "<", Pos: pos}
	}
}

// readStarStart handles tokens starting with '*'.
func (l *Lexer) readStarStart(pos Pos) Token {
	if l.isArrowContinuation(l.peekChar()) {
		return l.readArrowFrom(pos)
	}
	return l.readError(pos)
}

// readArrowOrDots handles sequences starting with '..' (dotted arrows or range).
func (l *Lexer) readArrowOrDots(pos Pos) Token {
	var b strings.Builder
	b.WriteRune(l.ch)
	l.readChar()
	for !l.eof && (l.ch == '.' || l.ch == '-') {
		b.WriteRune(l.ch)
		l.readChar()
	}
	// Check for arrowhead.
	if !l.eof {
		switch l.ch {
		case '>':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case '|':
			b.WriteRune(l.ch)
			l.readChar()
			if !l.eof && l.ch == '>' {
				b.WriteRune(l.ch)
				l.readChar()
			}
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case '*':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case 'o':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		}
	}
	return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
}

// continueArrow reads the shaft and optional arrowhead after a prefix.
func (l *Lexer) continueArrow(b *strings.Builder, pos Pos) Token {
	for !l.eof && (l.ch == '-' || l.ch == '.') {
		b.WriteRune(l.ch)
		l.readChar()
	}
	// Check for arrowhead at end.
	if !l.eof {
		switch l.ch {
		case '>':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case '|':
			b.WriteRune(l.ch)
			l.readChar()
			if !l.eof && l.ch == '>' {
				b.WriteRune(l.ch)
				l.readChar()
			}
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case '*':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		case 'o':
			b.WriteRune(l.ch)
			l.readChar()
			return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
		}
	}
	return Token{Type: TokenArrow, Literal: b.String(), Pos: pos}
}
