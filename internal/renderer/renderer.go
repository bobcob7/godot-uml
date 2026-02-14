// Package renderer defines the rendering interface for diagram output.
package renderer

import (
	"io"

	"github.com/bobcob7/godot-uml/internal/ast"
)

// Renderer produces output from a parsed diagram.
type Renderer interface {
	Render(w io.Writer, diagram *ast.Diagram) error
}
