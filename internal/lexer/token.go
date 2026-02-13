package lexer

import "fmt"

// TokenType classifies a lexical token.
type TokenType int

const (
	// Special tokens.
	TokenError TokenType = iota
	TokenEOF

	// Delimiters and punctuation.
	TokenLBrace    // {
	TokenRBrace    // }
	TokenLParen    // (
	TokenRParen    // )
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenColon     // :
	TokenComma     // ,
	TokenDot       // .
	TokenNewline   // \n
	TokenPipe      // |
	TokenHash      // # (also visibility)
	TokenLAngle    // <
	TokenRAngle    // >
	TokenEquals    // =
	TokenSemicolon // ;

	// Visibility markers.
	TokenPlus  // +
	TokenMinus // -
	TokenTilde // ~

	// Diagram delimiters.
	TokenStartUML // @startuml
	TokenEndUML   // @enduml

	// Class diagram keywords.
	TokenClass      // class
	TokenInterface  // interface
	TokenEnum       // enum
	TokenAbstract   // abstract
	TokenExtends    // extends
	TokenImplements // implements
	TokenPackage    // package
	TokenNamespace  // namespace
	TokenAs         // as
	TokenStatic     // {static}
	TokenField      // {field}
	TokenMethod     // {method}

	// Sequence diagram keywords.
	TokenParticipant // participant
	TokenActor       // actor
	TokenBoundary    // boundary
	TokenControl     // control
	TokenEntity      // entity
	TokenDatabase    // database
	TokenCollections // collections
	TokenQueue       // queue
	TokenActivate    // activate
	TokenDeactivate  // deactivate
	TokenReturn      // return
	TokenAlt         // alt
	TokenElse        // else
	TokenEnd         // end
	TokenLoop        // loop
	TokenGroup       // group
	TokenNote        // note
	TokenOf          // of
	TokenOver        // over
	TokenLeft        // left
	TokenRight       // right

	// Arrows.
	TokenArrow // ->, -->, <-, <--, <|--,  *--, o--, etc.

	// Directives.
	TokenSkinparam // skinparam
	TokenHide      // hide
	TokenShow      // show
	TokenTitle     // title
	TokenHeader    // header
	TokenFooter    // footer

	// Literals.
	TokenIdent  // identifiers
	TokenString // "..." or '...'
	TokenNumber // integer or decimal

	// Comments.
	TokenLineComment  // ' single-line comment
	TokenBlockComment // /' ... '/
)

//go:generate stringer -type=TokenType -trimprefix=Token

// Pos represents a source position.
type Pos struct {
	Line   int // 1-based line number
	Column int // 1-based column number
}

// String returns the position as "line:col".
func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Token is a lexical token with its type, literal text, and source position.
type Token struct {
	Type    TokenType
	Literal string
	Pos     Pos
}

// String returns a debug representation of the token.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%s", t.Type, t.Literal, t.Pos)
}
