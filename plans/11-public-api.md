COMPLETE # Plan 11 — Public Library API

**Status**: Complete

**Delivers**: The `pkg/goduml` public API that ties together parsing, validation, layout, and rendering into a simple interface for Go consumers.

## Acceptance Criteria

- [ ] `goduml.Render(io.Reader, io.Writer, ...Option) error` — primary API
- [ ] Options for: theme selection, output format (SVG), custom skinparams
- [ ] `goduml.Validate(io.Reader) []Error` — validate without rendering
- [ ] `goduml.Parse(io.Reader) (*ast.Diagram, []Error)` — parse and return AST
- [ ] Errors include source position and human-readable messages
- [ ] Integration tests exercise the full pipeline end-to-end
- [ ] Example usage in godoc

## Dependencies

- Plan 09
- Plan 10
