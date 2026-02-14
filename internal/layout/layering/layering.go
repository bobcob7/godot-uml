// Package layering implements layer assignment for the Sugiyama layout algorithm.
package layering

// Assign performs longest-path layer assignment on a DAG.
// adj is an adjacency list, n is the number of nodes.
// Returns a slice where layers[i] is the layer of node i (0-based).
func Assign(adj [][]int, n int) []int {
	layers := make([]int, n)
	inDeg := make([]int, n)
	for u := range n {
		for _, v := range adj[u] {
			inDeg[v]++
		}
	}
	queue := make([]int, 0, n)
	for i := range n {
		if inDeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	if len(queue) == 0 && n > 0 {
		queue = append(queue, 0)
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if layers[u]+1 > layers[v] {
				layers[v] = layers[u] + 1
			}
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return layers
}
