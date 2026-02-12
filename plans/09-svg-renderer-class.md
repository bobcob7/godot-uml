# Plan 09 — SVG Renderer for Class Diagrams

**Status**: Pending

**Delivers**: An SVG renderer that takes a positioned class diagram AST and produces valid SVG output with Darcula styling.

## Acceptance Criteria

- [ ] Renders classes as rounded rectangles with compartments (name, fields, methods)
- [ ] Renders visibility icons, stereotypes, and abstract/static modifiers
- [ ] Renders all relationship types with correct arrow heads (inheritance triangle, composition diamond, etc.)
- [ ] Renders relationship labels and cardinality
- [ ] Renders notes attached to elements
- [ ] Renders packages/namespaces as containing rectangles
- [ ] Output is valid SVG 1.1 viewable in browsers
- [ ] Darcula theme applied by default
- [ ] Integration tests compare output SVG against golden files
- [ ] Visual regression tests (SVG snapshot comparison)

## Dependencies

- Plan 04
- Plan 06
- Plan 07
- Plan 08
