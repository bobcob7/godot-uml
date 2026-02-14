package layering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	t.Parallel()
	t.Run("LinearChain", func(t *testing.T) {
		t.Parallel()
		// 0 -> 1 -> 2
		adj := [][]int{{1}, {2}, {}}
		layers := Assign(adj, 3)
		assert.Equal(t, 0, layers[0])
		assert.Equal(t, 1, layers[1])
		assert.Equal(t, 2, layers[2])
	})
	t.Run("Diamond", func(t *testing.T) {
		t.Parallel()
		// 0 -> 1, 0 -> 2, 1 -> 3, 2 -> 3
		adj := [][]int{{1, 2}, {3}, {3}, {}}
		layers := Assign(adj, 4)
		assert.Equal(t, 0, layers[0])
		assert.Equal(t, 1, layers[1])
		assert.Equal(t, 1, layers[2])
		assert.Equal(t, 2, layers[3])
	})
	t.Run("SingleNode", func(t *testing.T) {
		t.Parallel()
		adj := [][]int{{}}
		layers := Assign(adj, 1)
		assert.Equal(t, 0, layers[0])
	})
	t.Run("DisconnectedNodes", func(t *testing.T) {
		t.Parallel()
		adj := [][]int{{}, {}, {}}
		layers := Assign(adj, 3)
		assert.Equal(t, 0, layers[0])
		assert.Equal(t, 0, layers[1])
		assert.Equal(t, 0, layers[2])
	})
	t.Run("MultipleSources", func(t *testing.T) {
		t.Parallel()
		// 0 -> 2, 1 -> 2
		adj := [][]int{{2}, {2}, {}}
		layers := Assign(adj, 3)
		assert.Equal(t, 0, layers[0])
		assert.Equal(t, 0, layers[1])
		assert.Equal(t, 1, layers[2])
	})
	t.Run("LongPath", func(t *testing.T) {
		t.Parallel()
		// 0 -> 1 -> 2 -> 3 -> 4
		adj := [][]int{{1}, {2}, {3}, {4}, {}}
		layers := Assign(adj, 5)
		for i := range 5 {
			assert.Equal(t, i, layers[i])
		}
	})
	t.Run("EmptyGraph", func(t *testing.T) {
		t.Parallel()
		layers := Assign(nil, 0)
		assert.Empty(t, layers)
	})
}
