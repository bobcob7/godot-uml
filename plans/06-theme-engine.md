# Plan 06 — Theme Engine & Skinparam System

**Status**: Pending

**Delivers**: A theme/styling engine that resolves visual properties for any diagram element. Ships with the Darcula theme as default.

## Acceptance Criteria

- [ ] Theme defines colors, fonts, sizes, borders, spacing for all element types
- [ ] Skinparam directives in diagrams override theme defaults
- [ ] Property resolution: diagram skinparam → theme → hardcoded fallback
- [ ] Darcula theme matches JetBrains color palette
- [ ] Theme is loadable from a struct (library use) or from parsed skinparams
- [ ] Unit tests verify property resolution precedence
- [ ] Unit tests verify all Darcula palette values

## Darcula Palette

| Element | Color |
|---|---|
| Background | `#2B2B2B` |
| Text / Foreground | `#A9B7C6` |
| Class fill | `#3C3F41` |
| Border | `#555555` |
| Arrow / line | `#A9B7C6` |
| Note fill | `#4E5254` |
| Stereotype text | `#CC7832` |
| String / annotation | `#6A8759` |
| Interface text | `#6897BB` |

## Dependencies

- Plan 01
