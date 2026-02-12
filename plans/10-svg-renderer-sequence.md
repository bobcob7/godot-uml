# Plan 10 — SVG Renderer for Sequence Diagrams

**Status**: Pending

**Delivers**: An SVG renderer that takes a positioned sequence diagram AST and produces valid SVG output with Darcula styling.

## Acceptance Criteria

- [ ] Renders participant boxes at the top with lifeline dashes
- [ ] Renders messages as arrows between lifelines (solid for sync, dashed for async)
- [ ] Renders activation bars on lifelines
- [ ] Renders fragments (alt/loop/par) as labeled rectangles spanning participants
- [ ] Renders notes (left, right, over)
- [ ] Renders dividers, delays, and autonumber labels
- [ ] Sequence diagram uses vertical time axis layout (not Sugiyama — custom layout)
- [ ] Output is valid SVG 1.1
- [ ] Integration tests with golden SVG files

## Dependencies

- Plan 05
- Plan 06
- Plan 07
