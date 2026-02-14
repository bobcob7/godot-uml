# ✅ COMPLETE Plan 03 — AST Node Types & Parser Foundation

**Status**: Complete

**Delivers**: AST type definitions for all supported diagram elements, and a recursive descent parser framework that can parse `@startuml`/`@enduml` boundaries, comments, titles, and empty diagrams.

## Acceptance Criteria

- [ ] AST types cover: Diagram (root), ClassDef, InterfaceDef, EnumDef, Field, Method, Relationship, Participant, Message, Fragment (alt/loop/par), Note, Skinparam, Package/Namespace
- [ ] Parser correctly handles: empty diagrams, diagrams with only comments, title/header/footer directives
- [ ] Validation errors include source position and human-readable message
- [ ] Error recovery: parser continues after errors to report multiple issues
- [ ] Unit tests for each parser entry point

## Dependencies

- Plan 02
