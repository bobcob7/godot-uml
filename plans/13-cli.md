# Plan 13 — CLI Interface

**Status**: Pending

**Delivers**: A command-line interface for rendering `.puml` files to SVG, validating files, and starting the server.

## Acceptance Criteria

- [ ] `godot-uml render input.puml -o output.svg` — render a file
- [ ] `godot-uml render -` — read from stdin, write SVG to stdout
- [ ] `godot-uml validate input.puml` — validate and report errors
- [ ] `godot-uml serve --port 8080` — start the HTTP server
- [ ] `godot-uml version` — print version info
- [ ] Exit codes: 0 = success, 1 = validation errors, 2 = system error
- [ ] Help text for all commands
- [ ] Integration tests for CLI commands

## Dependencies

- Plan 11
- Plan 12
