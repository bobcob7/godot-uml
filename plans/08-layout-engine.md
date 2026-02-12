# Plan 08 — Sugiyama Layout Engine

**Status**: Pending

**Delivers**: A hierarchical graph layout algorithm (Sugiyama) that takes a graph of nodes with sizes and edges, and produces (x, y) positions for all elements.

## Acceptance Criteria

- [ ] Implements four Sugiyama phases: cycle removal, layer assignment, crossing minimization, coordinate assignment
- [ ] Handles: self-loops, multi-edges, disconnected subgraphs
- [ ] Nodes have configurable width/height (from text measurement)
- [ ] Edges have optional labels that affect spacing
- [ ] Produces non-overlapping layout with configurable padding/margins
- [ ] Unit tests for each phase independently
- [ ] Integration tests verify no overlapping nodes for various graph topologies

## Dependencies

- Plan 07
