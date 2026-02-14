// Package coordinate implements coordinate assignment for the Sugiyama layout algorithm.
package coordinate

// NodeSize holds the dimensions of a node.
type NodeSize struct {
	Width  float64
	Height float64
}

// Position holds the computed coordinates of a node.
type Position struct {
	X float64
	Y float64
}

// Assign computes X and Y positions for all nodes.
// layers contains node indices grouped by layer.
// sizes holds width/height for each node index.
// nodePadding is horizontal spacing between nodes.
// layerSpacing is vertical spacing between layers.
func Assign(layers [][]int, sizes []NodeSize, nodePadding, layerSpacing float64) []Position {
	n := len(sizes)
	positions := make([]Position, n)
	// First pass: assign positions left-to-right, top-to-bottom.
	y := 0.0
	for _, layer := range layers {
		x := 0.0
		maxHeight := 0.0
		for _, idx := range layer {
			positions[idx] = Position{X: x, Y: y}
			x += sizes[idx].Width + nodePadding
			if sizes[idx].Height > maxHeight {
				maxHeight = sizes[idx].Height
			}
		}
		y += maxHeight + layerSpacing
	}
	// Second pass: center layers relative to the widest.
	maxWidth := 0.0
	for _, layer := range layers {
		w := layerWidth(sizes, layer, nodePadding)
		if w > maxWidth {
			maxWidth = w
		}
	}
	for _, layer := range layers {
		w := layerWidth(sizes, layer, nodePadding)
		offset := (maxWidth - w) / 2
		for _, idx := range layer {
			positions[idx].X += offset
		}
	}
	return positions
}

func layerWidth(sizes []NodeSize, layer []int, padding float64) float64 {
	if len(layer) == 0 {
		return 0
	}
	w := 0.0
	for _, idx := range layer {
		w += sizes[idx].Width
	}
	w += padding * float64(len(layer)-1)
	return w
}
