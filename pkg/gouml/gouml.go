// Package gouml provides the public library API for rendering PlantUML diagrams to SVG.
//
// The primary entry point is Render, which reads PlantUML input and writes SVG output:
//
//	err := gouml.Render(os.Stdin, os.Stdout)
//
// Use options to customize rendering:
//
//	err := gouml.Render(input, output,
//	    gouml.WithSkinparam("backgroundColor", "#FFFFFF"),
//	)
//
// For validation without rendering:
//
//	errs := gouml.Validate(input)
//
// For parsing to an AST:
//
//	diagram, errs := gouml.Parse(input)
package gouml

import (
	"fmt"
	"io"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/bobcob7/go-uml/internal/parser"
	"github.com/bobcob7/go-uml/internal/renderer/svg"
	"github.com/bobcob7/go-uml/internal/theme"
)

// Diagram is an opaque handle to a parsed PlantUML diagram.
// Obtain one via Parse, then pass it to RenderDiagram.
type Diagram struct {
	internal *ast.Diagram
}

// Error represents a parse or validation error with source position.
type Error struct {
	Line    int
	Column  int
	Message string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Column, e.Message)
}

// Option configures rendering behavior.
type Option func(*options)

type options struct {
	theme      *theme.Theme
	skinparams map[string]string
}

// WithTheme sets the theme for rendering.
// If not specified, the Darcula theme is used.
func WithTheme(t *theme.Theme) Option {
	return func(o *options) {
		o.theme = t
	}
}

// WithSkinparam sets a skinparam override that takes highest priority
// in the property resolution chain.
func WithSkinparam(name, value string) Option {
	return func(o *options) {
		o.skinparams[name] = value
	}
}

// Render reads PlantUML from r and writes SVG to w.
// Options may be provided to customize theme and skinparam overrides.
func Render(r io.Reader, w io.Writer, opts ...Option) error {
	diagram, errs := Parse(r)
	if len(errs) > 0 {
		return errs[0]
	}
	return RenderDiagram(w, diagram, opts...)
}

// RenderDiagram renders a previously parsed diagram to SVG.
func RenderDiagram(w io.Writer, d *Diagram, opts ...Option) error {
	o := &options{skinparams: make(map[string]string)}
	for _, opt := range opts {
		opt(o)
	}
	resolver := theme.NewResolver(o.theme)
	for k, v := range o.skinparams {
		resolver.SetSkinparam(k, v)
	}
	if isSequenceDiagram(d.internal) {
		return svg.NewSequenceRenderer(resolver).Render(w, d.internal)
	}
	return svg.NewClassRenderer(resolver).Render(w, d.internal)
}

// Parse reads PlantUML from r and returns the parsed diagram and any errors.
// Parsing uses error recovery to continue after errors and report multiple issues.
func Parse(r io.Reader) (*Diagram, []*Error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, []*Error{{Line: 1, Column: 1, Message: fmt.Sprintf("reading input: %s", err)}}
	}
	diagram, parseErrs := parser.Parse(string(data))
	if len(parseErrs) > 0 {
		errs := make([]*Error, len(parseErrs))
		for i, pe := range parseErrs {
			errs[i] = &Error{
				Line:    pe.Pos.Line,
				Column:  pe.Pos.Column,
				Message: pe.Message,
			}
		}
		return &Diagram{internal: diagram}, errs
	}
	return &Diagram{internal: diagram}, nil
}

// Validate reads PlantUML from r and returns any parse errors without rendering.
func Validate(r io.Reader) []*Error {
	_, errs := Parse(r)
	return errs
}

// isSequenceDiagram inspects the AST to determine if it's a sequence diagram.
func isSequenceDiagram(d *ast.Diagram) bool {
	for _, stmt := range d.Statements {
		switch stmt.(type) {
		case *ast.Participant, *ast.Message, *ast.Fragment,
			*ast.Activate, *ast.Autonumber, *ast.Divider, *ast.Delay:
			return true
		}
	}
	return false
}
