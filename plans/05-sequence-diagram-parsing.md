# Plan 05 — Sequence Diagram Parsing

**Status**: Complete

**Delivers**: Full parsing of PlantUML sequence diagram syntax into the AST.

## Acceptance Criteria

- [ ] Parses participant types: `participant`, `actor`, `boundary`, `control`, `entity`, `database`, `collections`, `queue`
- [ ] Parses messages: `->`, `-->`, `->>`, `-->>`, `<-`, `<--` with labels
- [ ] Parses activation/deactivation, `activate`/`deactivate`, `++`/`--` shorthand
- [ ] Parses fragments: `alt`/`else`, `loop`, `par`, `break`, `group`, `ref`
- [ ] Parses `autonumber`, dividers (`==`), delays (`...`), notes
- [ ] Clear validation errors for: unclosed fragments, message to undeclared participant, mismatched activate/deactivate
- [ ] Integration tests using `.puml` fixture files

## Dependencies

- Plan 03
