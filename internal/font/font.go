// Package font provides font metrics and embedded TTF fonts for text measurement.
// Fonts are embedded in the binary via the golang.org/x/image/font/gofont packages,
// requiring no external font files at runtime.
package font

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Size represents pixel dimensions.
type Size struct {
	Width  float64
	Height float64
}

// Family identifies a font family.
type Family string

const (
	FamilySans Family = "sans"
	FamilyMono Family = "mono"
	FamilyBold Family = "bold"
)

// parsedFonts caches parsed opentype fonts.
var (
	parsedFontsMu sync.Mutex
	parsedFonts   = map[Family]*opentype.Font{}
)

func parsedFont(family Family) (*opentype.Font, error) {
	parsedFontsMu.Lock()
	defer parsedFontsMu.Unlock()
	if f, ok := parsedFonts[family]; ok {
		return f, nil
	}
	var data []byte
	switch family {
	case FamilyMono:
		data = gomono.TTF
	case FamilyBold:
		data = gobold.TTF
	default:
		data = goregular.TTF
	}
	f, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing font %s: %w", family, err)
	}
	parsedFonts[family] = f
	return f, nil
}

func newFace(family Family, fontSize float64) (font.Face, error) {
	f, err := parsedFont(family)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("creating face for %s at %.1f: %w", family, fontSize, err)
	}
	return face, nil
}

// MeasureText computes the pixel dimensions for the given text at the specified
// font size and family. Multi-line text (separated by \n) is handled by measuring
// each line independently and returning the maximum width and total height.
func MeasureText(text string, fontSize float64, family Family) (Size, error) {
	face, err := newFace(family, fontSize)
	if err != nil {
		return Size{}, err
	}
	defer func() { _ = face.Close() }()
	metrics := face.Metrics()
	lineHeight := fixedToFloat(metrics.Height)
	if text == "" {
		return Size{Width: 0, Height: lineHeight}, nil
	}
	lines := strings.Split(text, "\n")
	var maxWidth float64
	for _, line := range lines {
		w := fixedToFloat(font.MeasureString(face, line))
		if w > maxWidth {
			maxWidth = w
		}
	}
	totalHeight := float64(len(lines)) * lineHeight
	return Size{Width: maxWidth, Height: totalHeight}, nil
}

func fixedToFloat(v fixed.Int26_6) float64 {
	return float64(v) / 64.0
}
