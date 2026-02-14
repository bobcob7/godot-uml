// Package crossing implements crossing minimization for the Sugiyama layout algorithm.
package crossing

// Count returns the number of edge crossings between two adjacent layers.
// upper and lower are node indices in their respective layers.
// adj maps upper-layer nodes to their lower-layer neighbors.
func Count(upper, lower []int, adj map[int][]int) int {
	// Build position map for lower layer.
	pos := make(map[int]int, len(lower))
	for i, node := range lower {
		pos[node] = i
	}
	// Collect edges as (upper_pos, lower_pos) pairs.
	type pair struct{ u, v int }
	var edges []pair
	for i, node := range upper {
		for _, neighbor := range adj[node] {
			if p, ok := pos[neighbor]; ok {
				edges = append(edges, pair{i, p})
			}
		}
	}
	// Count inversions (crossings).
	crossings := 0
	for i := range edges {
		for j := i + 1; j < len(edges); j++ {
			if (edges[i].u < edges[j].u && edges[i].v > edges[j].v) ||
				(edges[i].u > edges[j].u && edges[i].v < edges[j].v) {
				crossings++
			}
		}
	}
	return crossings
}

// Minimize reorders nodes in each layer to reduce edge crossings
// using the barycenter heuristic with multiple passes.
// layers contains node indices grouped by layer.
// adj is the full adjacency list (node → neighbors).
// n is the total number of nodes.
// Returns the reordered layers.
func Minimize(layers [][]int, adj [][]int, n int) [][]int {
	radj := buildReverse(adj, n)
	result := make([][]int, len(layers))
	for i := range layers {
		result[i] = append([]int(nil), layers[i]...)
	}
	for pass := range 4 {
		if pass%2 == 0 {
			for i := 0; i < len(result)-1; i++ {
				result[i+1] = reorder(result[i+1], result[i], radj)
			}
		} else {
			for i := len(result) - 1; i > 0; i-- {
				result[i-1] = reorder(result[i-1], result[i], adj)
			}
		}
	}
	return result
}

func buildReverse(adj [][]int, n int) [][]int {
	radj := make([][]int, n)
	for u, neighbors := range adj {
		for _, v := range neighbors {
			if v < n {
				radj[v] = append(radj[v], u)
			}
		}
	}
	return radj
}

func reorder(layer, fixedLayer []int, connections [][]int) []int {
	pos := make(map[int]int, len(fixedLayer))
	for i, node := range fixedLayer {
		pos[node] = i
	}
	type entry struct {
		node       int
		barycenter float64
	}
	entries := make([]entry, 0, len(layer))
	for _, node := range layer {
		sum := 0.0
		count := 0
		if node < len(connections) {
			for _, neighbor := range connections[node] {
				if p, ok := pos[neighbor]; ok {
					sum += float64(p)
					count++
				}
			}
		}
		bc := float64(len(entries))
		if count > 0 {
			bc = sum / float64(count)
		}
		entries = append(entries, entry{node: node, barycenter: bc})
	}
	// Insertion sort for stability.
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].barycenter < entries[j-1].barycenter; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
	result := make([]int, len(entries))
	for i, e := range entries {
		result[i] = e.node
	}
	return result
}
