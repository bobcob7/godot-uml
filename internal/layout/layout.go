// Package layout implements the Sugiyama hierarchical layout algorithm for graph positioning.
package layout

// Node represents a graph node with dimensions.
type Node struct {
	ID      string
	Width   float64
	Height  float64
	Virtual bool // true for virtual nodes inserted for long edges
	X       float64
	Y       float64
	Layer   int
	Order   int
}

// Edge represents a directed edge between two nodes.
type Edge struct {
	From     string
	To       string
	Label    string
	Reversed bool // true if edge was reversed during cycle removal
}

// Graph represents the input graph for layout.
type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

// Options configures the layout algorithm.
type Options struct {
	NodePadding  float64 // horizontal spacing between nodes in a layer
	LayerSpacing float64 // vertical spacing between layers
}

// DefaultOptions returns sensible default layout options.
func DefaultOptions() Options {
	return Options{
		NodePadding:  40,
		LayerSpacing: 60,
	}
}

// Layout runs the full Sugiyama algorithm on the graph.
// It modifies the nodes in place, setting their X, Y, Layer, and Order fields.
func Layout(g *Graph, opts Options) {
	if len(g.Nodes) == 0 {
		return
	}
	nodeIndex := buildNodeIndex(g)
	adj := buildAdjacency(g, nodeIndex)
	n := len(g.Nodes)
	// Phase 1: Cycle removal.
	reversed := removeCycles(adj, n)
	for _, e := range g.Edges {
		key := edgeKey(nodeIndex[e.From], nodeIndex[e.To])
		if reversed[key] {
			e.Reversed = true
		}
	}
	// Phase 2: Layer assignment.
	layers := assignLayers(adj, n)
	for i, layer := range layers {
		g.Nodes[i].Layer = layer
	}
	// Phase 3: Insert virtual nodes for long edges.
	adj, layers, g.Nodes = insertVirtualNodes(adj, layers, g.Nodes, opts)
	// Phase 4: Order nodes within layers.
	layerBuckets := buildLayerBuckets(layers)
	layerBuckets = minimizeCrossings(layerBuckets, adj, len(g.Nodes))
	for order, idx := range flattenBuckets(layerBuckets) {
		_ = order
		g.Nodes[idx].Order = orderInLayer(layerBuckets, layers[idx], idx)
	}
	// Phase 5: Coordinate assignment.
	assignCoordinates(g.Nodes, layerBuckets, opts)
}

func buildNodeIndex(g *Graph) map[string]int {
	idx := make(map[string]int, len(g.Nodes))
	for i, n := range g.Nodes {
		idx[n.ID] = i
	}
	return idx
}

func buildAdjacency(g *Graph, nodeIndex map[string]int) [][]int {
	n := len(g.Nodes)
	adj := make([][]int, n)
	for _, e := range g.Edges {
		from, okF := nodeIndex[e.From]
		to, okT := nodeIndex[e.To]
		if !okF || !okT {
			continue
		}
		if from == to {
			continue // skip self-loops
		}
		adj[from] = append(adj[from], to)
	}
	return adj
}

func edgeKey(from, to int) [2]int {
	return [2]int{from, to}
}

// removeCycles uses DFS to find back edges and reverses them.
// Returns a set of reversed edge keys.
func removeCycles(adj [][]int, n int) map[[2]int]bool {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make([]int, n)
	reversed := make(map[[2]int]bool)
	var dfs func(u int)
	dfs = func(u int) {
		color[u] = gray
		newAdj := make([]int, 0, len(adj[u]))
		for _, v := range adj[u] {
			switch color[v] {
			case gray:
				// Back edge — reverse it.
				reversed[edgeKey(u, v)] = true
				adj[v] = append(adj[v], u)
			case white:
				newAdj = append(newAdj, v)
				dfs(v)
			default:
				newAdj = append(newAdj, v)
			}
		}
		adj[u] = newAdj
		color[u] = black
	}
	for i := range n {
		if color[i] == white {
			dfs(i)
		}
	}
	return reversed
}

// assignLayers uses the longest path algorithm from sources.
func assignLayers(adj [][]int, n int) []int {
	layers := make([]int, n)
	// Compute in-degrees.
	inDeg := make([]int, n)
	for u := range n {
		for _, v := range adj[u] {
			inDeg[v]++
		}
	}
	// BFS from sources (in-degree 0).
	queue := make([]int, 0, n)
	for i := range n {
		if inDeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	// Handle disconnected graphs with no sources.
	if len(queue) == 0 {
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

// insertVirtualNodes adds dummy nodes for edges spanning more than one layer.
func insertVirtualNodes(adj [][]int, layers []int, nodes []*Node, opts Options) ([][]int, []int, []*Node) {
	newAdj := make([][]int, len(adj))
	for i := range adj {
		newAdj[i] = append([]int(nil), adj[i]...)
	}
	for u := range len(adj) {
		for j, v := range adj[u] {
			span := layers[v] - layers[u]
			if span <= 1 {
				continue
			}
			// Replace long edge with chain of virtual nodes.
			prev := u
			for k := 1; k < span; k++ {
				vn := &Node{
					ID:      "",
					Width:   0,
					Height:  0,
					Virtual: true,
					Layer:   layers[u] + k,
				}
				vnIdx := len(nodes)
				nodes = append(nodes, vn)
				layers = append(layers, layers[u]+k)
				newAdj = append(newAdj, nil)
				// Remove old edge from prev to v.
				if prev == u {
					newAdj[prev][j] = vnIdx
				} else {
					newAdj[prev] = append(newAdj[prev], vnIdx)
				}
				prev = vnIdx
			}
			newAdj[prev] = append(newAdj[prev], v)
		}
	}
	_ = opts
	return newAdj, layers, nodes
}

func buildLayerBuckets(layers []int) [][]int {
	maxLayer := 0
	for _, l := range layers {
		if l > maxLayer {
			maxLayer = l
		}
	}
	buckets := make([][]int, maxLayer+1)
	for i, l := range layers {
		buckets[l] = append(buckets[l], i)
	}
	return buckets
}

// minimizeCrossings applies the barycenter heuristic.
func minimizeCrossings(layerBuckets [][]int, adj [][]int, n int) [][]int {
	// Build reverse adjacency.
	radj := make([][]int, n)
	for u, neighbors := range adj {
		for _, v := range neighbors {
			if v < n {
				radj[v] = append(radj[v], u)
			}
		}
	}
	// Grow radj/adj if needed for virtual nodes.
	for len(radj) < n {
		radj = append(radj, nil)
	}
	// Multiple passes, alternating down and up sweeps.
	for pass := range 4 {
		if pass%2 == 0 {
			// Down sweep: fix layer i, reorder layer i+1.
			for i := 0; i < len(layerBuckets)-1; i++ {
				layerBuckets[i+1] = reorderByBarycenter(layerBuckets[i+1], layerBuckets[i], radj)
			}
		} else {
			// Up sweep: fix layer i, reorder layer i-1.
			for i := len(layerBuckets) - 1; i > 0; i-- {
				layerBuckets[i-1] = reorderByBarycenter(layerBuckets[i-1], layerBuckets[i], adj)
			}
		}
	}
	return layerBuckets
}

func reorderByBarycenter(layer, fixedLayer []int, connections [][]int) []int {
	// Build position map for fixed layer.
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
		for _, neighbor := range connections[node] {
			if p, ok := pos[neighbor]; ok {
				sum += float64(p)
				count++
			}
		}
		bc := float64(len(entries)) // default: keep position
		if count > 0 {
			bc = sum / float64(count)
		}
		entries = append(entries, entry{node: node, barycenter: bc})
	}
	// Sort by barycenter (stable to preserve relative order for ties).
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

func flattenBuckets(buckets [][]int) []int {
	var all []int
	for _, b := range buckets {
		all = append(all, b...)
	}
	return all
}

func orderInLayer(buckets [][]int, layer, nodeIdx int) int {
	for i, idx := range buckets[layer] {
		if idx == nodeIdx {
			return i
		}
	}
	return 0
}

// assignCoordinates sets X and Y positions for all nodes.
func assignCoordinates(nodes []*Node, layerBuckets [][]int, opts Options) {
	y := 0.0
	for _, layer := range layerBuckets {
		x := 0.0
		maxHeight := 0.0
		for _, idx := range layer {
			node := nodes[idx]
			node.X = x
			node.Y = y
			x += node.Width + opts.NodePadding
			if node.Height > maxHeight {
				maxHeight = node.Height
			}
		}
		y += maxHeight + opts.LayerSpacing
	}
	// Center layers horizontally relative to the widest layer.
	maxWidth := 0.0
	for _, layer := range layerBuckets {
		w := layerWidth(nodes, layer, opts.NodePadding)
		if w > maxWidth {
			maxWidth = w
		}
	}
	for _, layer := range layerBuckets {
		w := layerWidth(nodes, layer, opts.NodePadding)
		offset := (maxWidth - w) / 2
		for _, idx := range layer {
			nodes[idx].X += offset
		}
	}
}

func layerWidth(nodes []*Node, layer []int, padding float64) float64 {
	if len(layer) == 0 {
		return 0
	}
	w := 0.0
	for _, idx := range layer {
		w += nodes[idx].Width
	}
	w += padding * float64(len(layer)-1)
	return w
}
