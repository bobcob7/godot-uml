# Plan 12 — HTTP Server & Live Editor

**Status**: Pending

**Delivers**: A standalone HTTP server with PlantUML-compatible API endpoints and a minimal live editor (textarea + SVG preview).

## Acceptance Criteria

- [ ] `POST /render` accepts PlantUML text body, returns SVG
- [ ] `GET /svg/{encoded}` accepts PlantUML-encoded URL, returns SVG
- [ ] `GET /` serves the live editor HTML page
- [ ] Editor: textarea on the left, live SVG preview on the right, auto-renders on keystroke (debounced)
- [ ] PlantUML URL encoding/decoding (DEFLATE + custom base64 alphabet)
- [ ] Validation errors displayed in the editor with line numbers
- [ ] Static assets embedded in the binary (`embed` package)
- [ ] Server configurable: port, host, read/write timeouts
- [ ] Integration tests for all HTTP endpoints
- [ ] The binary is a single self-contained executable

## Dependencies

- Plan 11
