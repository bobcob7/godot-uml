package coordinate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	t.Parallel()
	t.Run("SingleNode", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{{Width: 100, Height: 50}}
		layers := [][]int{{0}}
		positions := Assign(layers, sizes, 40, 60)
		assert.Equal(t, 0.0, positions[0].X)
		assert.Equal(t, 0.0, positions[0].Y)
	})
	t.Run("TwoNodesOneLayer", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{{Width: 100, Height: 50}, {Width: 100, Height: 50}}
		layers := [][]int{{0, 1}}
		positions := Assign(layers, sizes, 40, 60)
		assert.Equal(t, positions[0].Y, positions[1].Y, "same layer = same Y")
		assert.Greater(t, positions[1].X, positions[0].X, "second node to the right")
	})
	t.Run("TwoLayers", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{{Width: 100, Height: 50}, {Width: 100, Height: 50}}
		layers := [][]int{{0}, {1}}
		positions := Assign(layers, sizes, 40, 60)
		assert.Equal(t, positions[0].X, positions[1].X, "single-node layers centered equally")
		assert.Greater(t, positions[1].Y, positions[0].Y, "second layer below first")
		assert.Equal(t, 110.0, positions[1].Y, "Y = height(50) + spacing(60)")
	})
	t.Run("CentersNarrowerLayers", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{
			{Width: 100, Height: 50},
			{Width: 100, Height: 50},
			{Width: 50, Height: 30},
		}
		layers := [][]int{{0, 1}, {2}}
		positions := Assign(layers, sizes, 40, 60)
		// Top layer width: 100 + 40 + 100 = 240
		// Bottom layer width: 50
		// Offset: (240 - 50) / 2 = 95
		assert.Equal(t, 95.0, positions[2].X, "narrow layer centered under wide layer")
	})
	t.Run("NoOverlap", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{
			{Width: 80, Height: 40},
			{Width: 120, Height: 60},
			{Width: 100, Height: 50},
		}
		layers := [][]int{{0, 1, 2}}
		positions := Assign(layers, sizes, 20, 40)
		// Node 0 ends at X + Width.
		assert.GreaterOrEqual(t, positions[1].X, positions[0].X+sizes[0].Width, "no horizontal overlap between node 0 and 1")
		assert.GreaterOrEqual(t, positions[2].X, positions[1].X+sizes[1].Width, "no horizontal overlap between node 1 and 2")
	})
	t.Run("EmptyGraph", func(t *testing.T) {
		t.Parallel()
		positions := Assign(nil, nil, 40, 60)
		assert.Empty(t, positions)
	})
	t.Run("VariableHeightsAffectSpacing", func(t *testing.T) {
		t.Parallel()
		sizes := []NodeSize{
			{Width: 100, Height: 30},
			{Width: 100, Height: 80},
			{Width: 100, Height: 50},
		}
		layers := [][]int{{0, 1}, {2}}
		positions := Assign(layers, sizes, 40, 60)
		// Max height of layer 0 is 80, so layer 1 Y = 80 + 60 = 140.
		assert.Equal(t, 140.0, positions[2].Y, "Y accounts for tallest node in previous layer")
	})
}
