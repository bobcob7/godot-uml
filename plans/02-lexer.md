# Plan 02 — Lexer & Token Types

**Status**: Pending

**Delivers**: A lexer that tokenizes PlantUML text into a stream of typed tokens with accurate source positions (line, column). Supports both class diagram and sequence diagram syntax tokens.

## Acceptance Criteria

- [ ] Tokenizes: keywords (`@startuml`, `@enduml`, `class`, `interface`, `enum`, `abstract`, `participant`, `actor`, `->`, `-->`, `<|--`, `*--`, `o--`, etc.), identifiers, strings, braces, colons, visibility markers (`+`, `-`, `#`, `~`), comments, skinparam, directives
- [ ] Each token carries source position (line:col) for error reporting
- [ ] Unknown characters produce clear error tokens (not panics)
- [ ] Comprehensive table-driven unit tests covering all token types
- [ ] Tests for malformed input producing meaningful error tokens

## Dependencies

- Plan 01
