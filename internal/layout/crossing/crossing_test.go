package crossing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	t.Parallel()
	t.Run("NoCrossings", func(t *testing.T) {
		t.Parallel()
		// 0->2, 1->3 — parallel edges, no crossing.
		upper := []int{0, 1}
		lower := []int{2, 3}
		adj := map[int][]int{0: {2}, 1: {3}}
		assert.Equal(t, 0, Count(upper, lower, adj))
	})
	t.Run("OneCrossing", func(t *testing.T) {
		t.Parallel()
		// 0->3, 1->2 — crossed edges.
		upper := []int{0, 1}
		lower := []int{2, 3}
		adj := map[int][]int{0: {3}, 1: {2}}
		assert.Equal(t, 1, Count(upper, lower, adj))
	})
	t.Run("MultipleCrossings", func(t *testing.T) {
		t.Parallel()
		// 0->4, 1->3, 2->2 — reversed order.
		upper := []int{0, 1, 2}
		lower := []int{2, 3, 4}
		adj := map[int][]int{0: {4}, 1: {3}, 2: {2}}
		assert.Equal(t, 3, Count(upper, lower, adj))
	})
	t.Run("EmptyLayers", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 0, Count(nil, nil, nil))
	})
	t.Run("NoEdges", func(t *testing.T) {
		t.Parallel()
		upper := []int{0, 1}
		lower := []int{2, 3}
		adj := map[int][]int{}
		assert.Equal(t, 0, Count(upper, lower, adj))
	})
}

func TestMinimize(t *testing.T) {
	t.Parallel()
	t.Run("ReducesCrossings", func(t *testing.T) {
		t.Parallel()
		// Two layers: [0,1] -> [2,3], with crossed edges.
		adj := [][]int{{3}, {2}, {}, {}}
		layers := [][]int{{0, 1}, {2, 3}}
		beforeAdj := map[int][]int{0: {3}, 1: {2}}
		before := Count(layers[0], layers[1], beforeAdj)
		result := Minimize(layers, adj, 4)
		after := Count(result[0], result[1], beforeAdj)
		assert.LessOrEqual(t, after, before)
	})
	t.Run("PreservesAllNodes", func(t *testing.T) {
		t.Parallel()
		adj := [][]int{{2}, {3}, {}, {}}
		layers := [][]int{{0, 1}, {2, 3}}
		result := Minimize(layers, adj, 4)
		allNodes := make(map[int]bool)
		for _, layer := range result {
			for _, node := range layer {
				allNodes[node] = true
			}
		}
		assert.Len(t, allNodes, 4)
	})
	t.Run("SingleLayer", func(t *testing.T) {
		t.Parallel()
		adj := [][]int{{}, {}, {}}
		layers := [][]int{{0, 1, 2}}
		result := Minimize(layers, adj, 3)
		assert.Len(t, result, 1)
		assert.Len(t, result[0], 3)
	})
}
