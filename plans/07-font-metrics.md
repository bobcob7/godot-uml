# Plan 07 — Font Metrics & Text Measurement

**Status**: Pending

**Delivers**: Embedded font(s) and a text measurement service that computes pixel dimensions for text at any font size. Essential for layout accuracy.

## Acceptance Criteria

- [ ] At least one TTF font embedded in the binary (e.g., DejaVu Sans Mono or similar permissive-license font)
- [ ] `MeasureText(text, fontSize, fontFamily) → (width, height)` returns accurate pixel dimensions
- [ ] Handles multi-line text measurement
- [ ] Unit tests verify known text dimensions
- [ ] No external font files required at runtime

## Dependencies

- Plan 01
