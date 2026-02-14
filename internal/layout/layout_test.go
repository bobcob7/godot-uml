package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	t.Parallel()
	t.Run("EmptyGraph", func(t *testing.T) {
		t.Parallel()
		g := &Graph{}
		Layout(g, DefaultOptions())
		assert.Empty(t, g.Nodes)
	})
	t.Run("SingleNode", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{{ID: "A", Width: 100, Height: 50}},
		}
		Layout(g, DefaultOptions())
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 0.0, g.Nodes[0].X)
		assert.Equal(t, 0.0, g.Nodes[0].Y)
	})
	t.Run("LinearChain", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
				{ID: "C", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "B", To: "C"},
			},
		}
		Layout(g, DefaultOptions())
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 1, g.Nodes[1].Layer)
		assert.Equal(t, 2, g.Nodes[2].Layer)
		assert.Less(t, g.Nodes[0].Y, g.Nodes[1].Y)
		assert.Less(t, g.Nodes[1].Y, g.Nodes[2].Y)
	})
	t.Run("DiamondGraph", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
				{ID: "C", Width: 100, Height: 50},
				{ID: "D", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "A", To: "C"},
				{From: "B", To: "D"},
				{From: "C", To: "D"},
			},
		}
		Layout(g, DefaultOptions())
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 1, g.Nodes[1].Layer)
		assert.Equal(t, 1, g.Nodes[2].Layer)
		assert.Equal(t, 2, g.Nodes[3].Layer)
	})
	t.Run("CyclicGraph", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
				{ID: "C", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "B", To: "C"},
				{From: "C", To: "A"}, // cycle
			},
		}
		Layout(g, DefaultOptions())
		// Should handle cycle without infinite loop.
		// At least one edge should be reversed.
		hasReversed := false
		for _, e := range g.Edges {
			if e.Reversed {
				hasReversed = true
			}
		}
		assert.True(t, hasReversed, "should reverse at least one edge to break cycle")
	})
	t.Run("SelfLoop", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "A"}, // self-loop
				{From: "A", To: "B"},
			},
		}
		Layout(g, DefaultOptions())
		// Should not crash on self-loops.
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 1, g.Nodes[1].Layer)
	})
	t.Run("DisconnectedSubgraphs", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
				{ID: "C", Width: 100, Height: 50},
				{ID: "D", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "C", To: "D"},
			},
		}
		Layout(g, DefaultOptions())
		// Both subgraphs should be laid out.
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 1, g.Nodes[1].Layer)
		assert.Equal(t, 0, g.Nodes[2].Layer)
		assert.Equal(t, 1, g.Nodes[3].Layer)
	})
	t.Run("MultiEdges", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "A", To: "B"},
			},
		}
		Layout(g, DefaultOptions())
		assert.Equal(t, 0, g.Nodes[0].Layer)
		assert.Equal(t, 1, g.Nodes[1].Layer)
	})
}

func TestNoOverlap(t *testing.T) {
	t.Parallel()
	t.Run("WideGraph", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 150, Height: 60},
				{ID: "B", Width: 120, Height: 40},
				{ID: "C", Width: 180, Height: 50},
				{ID: "D", Width: 100, Height: 45},
				{ID: "E", Width: 130, Height: 55},
			},
			Edges: []*Edge{
				{From: "A", To: "B"},
				{From: "A", To: "C"},
				{From: "B", To: "D"},
				{From: "C", To: "D"},
				{From: "D", To: "E"},
			},
		}
		Layout(g, DefaultOptions())
		assertNoOverlap(t, g)
	})
	t.Run("FlatGraph", func(t *testing.T) {
		t.Parallel()
		g := &Graph{
			Nodes: []*Node{
				{ID: "A", Width: 100, Height: 50},
				{ID: "B", Width: 100, Height: 50},
				{ID: "C", Width: 100, Height: 50},
				{ID: "D", Width: 100, Height: 50},
			},
		}
		Layout(g, DefaultOptions())
		assertNoOverlap(t, g)
	})
	t.Run("LongChain", func(t *testing.T) {
		t.Parallel()
		g := &Graph{}
		for i := range 8 {
			g.Nodes = append(g.Nodes, &Node{
				ID:     string(rune('A' + i)),
				Width:  80 + float64(i*10),
				Height: 40 + float64(i*5),
			})
			if i > 0 {
				g.Edges = append(g.Edges, &Edge{
					From: string(rune('A' + i - 1)),
					To:   string(rune('A' + i)),
				})
			}
		}
		Layout(g, DefaultOptions())
		assertNoOverlap(t, g)
	})
}

func TestDefaultOptions(t *testing.T) {
	t.Parallel()
	opts := DefaultOptions()
	assert.Greater(t, opts.NodePadding, 0.0)
	assert.Greater(t, opts.LayerSpacing, 0.0)
}

func assertNoOverlap(t *testing.T, g *Graph) {
	t.Helper()
	realNodes := make([]*Node, 0)
	for _, n := range g.Nodes {
		if !n.Virtual {
			realNodes = append(realNodes, n)
		}
	}
	for i := range realNodes {
		for j := i + 1; j < len(realNodes); j++ {
			a := realNodes[i]
			b := realNodes[j]
			if a.Layer != b.Layer {
				continue
			}
			overlapX := a.X < b.X+b.Width && b.X < a.X+a.Width
			overlapY := a.Y < b.Y+b.Height && b.Y < a.Y+a.Height
			require.False(t, overlapX && overlapY,
				"nodes %s and %s overlap: (%v,%v,%v,%v) vs (%v,%v,%v,%v)",
				a.ID, b.ID, a.X, a.Y, a.Width, a.Height, b.X, b.Y, b.Width, b.Height)
		}
	}
}
