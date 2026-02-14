// Package theme provides theme definitions and skinparam resolution for diagram styling.
package theme

// Theme defines the complete visual styling for diagram rendering.
// Property resolution order: skinparam overrides → theme values → hardcoded fallbacks.
type Theme struct {
	// Global
	BackgroundColor string
	FontName        string
	FontSize        int
	FontColor       string
	// Class elements
	ClassBackgroundColor     string
	ClassBorderColor         string
	ClassFontColor           string
	ClassFontSize            int
	ClassStereotypeFontColor string
	// Interface elements
	InterfaceBackgroundColor string
	InterfaceBorderColor     string
	InterfaceFontColor       string
	// Enum elements
	EnumBackgroundColor string
	EnumBorderColor     string
	EnumFontColor       string
	// Arrows and lines
	ArrowColor    string
	ArrowFontSize int
	// Notes
	NoteBackgroundColor string
	NoteBorderColor     string
	NoteFontColor       string
	// Sequence diagram
	ParticipantBackgroundColor  string
	ParticipantBorderColor      string
	ParticipantFontColor        string
	SequenceLifeLineBorderColor string
	// Package/Namespace
	PackageBackgroundColor string
	PackageBorderColor     string
	PackageFontColor       string
	// Spacing
	Padding        int
	ClassPadding   int
	NotePadding    int
	BorderWidth    int
	ArrowThickness int
	// Annotation/string colors
	AnnotationColor string
}

// Darcula returns the default Darcula theme matching JetBrains color palette.
func Darcula() *Theme {
	return &Theme{
		BackgroundColor:             "#2B2B2B",
		FontName:                    "DejaVu Sans",
		FontSize:                    13,
		FontColor:                   "#A9B7C6",
		ClassBackgroundColor:        "#3C3F41",
		ClassBorderColor:            "#555555",
		ClassFontColor:              "#A9B7C6",
		ClassFontSize:               13,
		ClassStereotypeFontColor:    "#CC7832",
		InterfaceBackgroundColor:    "#3C3F41",
		InterfaceBorderColor:        "#555555",
		InterfaceFontColor:          "#6897BB",
		EnumBackgroundColor:         "#3C3F41",
		EnumBorderColor:             "#555555",
		EnumFontColor:               "#A9B7C6",
		ArrowColor:                  "#A9B7C6",
		ArrowFontSize:               11,
		NoteBackgroundColor:         "#4E5254",
		NoteBorderColor:             "#555555",
		NoteFontColor:               "#A9B7C6",
		ParticipantBackgroundColor:  "#3C3F41",
		ParticipantBorderColor:      "#555555",
		ParticipantFontColor:        "#A9B7C6",
		SequenceLifeLineBorderColor: "#555555",
		PackageBackgroundColor:      "#2B2B2B",
		PackageBorderColor:          "#555555",
		PackageFontColor:            "#A9B7C6",
		Padding:                     10,
		ClassPadding:                8,
		NotePadding:                 8,
		BorderWidth:                 1,
		ArrowThickness:              1,
		AnnotationColor:             "#6A8759",
	}
}

// hardcodedFallback returns the minimal fallback theme used when no theme is set.
func hardcodedFallback() *Theme {
	return &Theme{
		BackgroundColor:             "#FFFFFF",
		FontName:                    "sans-serif",
		FontSize:                    12,
		FontColor:                   "#000000",
		ClassBackgroundColor:        "#FEFECE",
		ClassBorderColor:            "#A80036",
		ClassFontColor:              "#000000",
		ClassFontSize:               12,
		ClassStereotypeFontColor:    "#000000",
		InterfaceBackgroundColor:    "#FEFECE",
		InterfaceBorderColor:        "#A80036",
		InterfaceFontColor:          "#000000",
		EnumBackgroundColor:         "#FEFECE",
		EnumBorderColor:             "#A80036",
		EnumFontColor:               "#000000",
		ArrowColor:                  "#000000",
		ArrowFontSize:               11,
		NoteBackgroundColor:         "#FBFB77",
		NoteBorderColor:             "#A80036",
		NoteFontColor:               "#000000",
		ParticipantBackgroundColor:  "#FEFECE",
		ParticipantBorderColor:      "#A80036",
		ParticipantFontColor:        "#000000",
		SequenceLifeLineBorderColor: "#A80036",
		PackageBackgroundColor:      "#FFFFFF",
		PackageBorderColor:          "#000000",
		PackageFontColor:            "#000000",
		Padding:                     10,
		ClassPadding:                8,
		NotePadding:                 8,
		BorderWidth:                 1,
		ArrowThickness:              1,
		AnnotationColor:             "#000000",
	}
}

// Resolver resolves style properties using a three-level hierarchy:
// skinparam overrides → theme → hardcoded fallback.
type Resolver struct {
	theme      *Theme
	fallback   *Theme
	skinparams map[string]string
}

// NewResolver creates a Resolver with the given theme.
// If theme is nil, the Darcula theme is used as default.
func NewResolver(theme *Theme) *Resolver {
	if theme == nil {
		theme = Darcula()
	}
	return &Resolver{
		theme:      theme,
		fallback:   hardcodedFallback(),
		skinparams: make(map[string]string),
	}
}

// SetSkinparam sets a skinparam override that takes highest priority.
func (r *Resolver) SetSkinparam(name, value string) {
	r.skinparams[name] = value
}

// skinparamKeys maps Theme field purpose to the skinparam name PlantUML uses.
var skinparamKeys = map[string]string{
	"BackgroundColor":             "backgroundColor",
	"FontName":                    "defaultFontName",
	"FontSize":                    "defaultFontSize",
	"FontColor":                   "defaultFontColor",
	"ClassBackgroundColor":        "classBackgroundColor",
	"ClassBorderColor":            "classBorderColor",
	"ClassFontColor":              "classFontColor",
	"ClassFontSize":               "classFontSize",
	"ClassStereotypeFontColor":    "classStereotypeFontColor",
	"InterfaceBackgroundColor":    "interfaceBackgroundColor",
	"InterfaceBorderColor":        "interfaceBorderColor",
	"InterfaceFontColor":          "interfaceFontColor",
	"EnumBackgroundColor":         "enumBackgroundColor",
	"EnumBorderColor":             "enumBorderColor",
	"EnumFontColor":               "enumFontColor",
	"ArrowColor":                  "arrowColor",
	"ArrowFontSize":               "arrowFontSize",
	"NoteBackgroundColor":         "noteBackgroundColor",
	"NoteBorderColor":             "noteBorderColor",
	"NoteFontColor":               "noteFontColor",
	"ParticipantBackgroundColor":  "participantBackgroundColor",
	"ParticipantBorderColor":      "participantBorderColor",
	"ParticipantFontColor":        "participantFontColor",
	"SequenceLifeLineBorderColor": "sequenceLifeLineBorderColor",
	"PackageBackgroundColor":      "packageBackgroundColor",
	"PackageBorderColor":          "packageBorderColor",
	"PackageFontColor":            "packageFontColor",
	"AnnotationColor":             "annotationColor",
}

// ResolveColor returns the color for a named property.
// Resolution order: skinparam → theme → fallback.
func (r *Resolver) ResolveColor(property string) string {
	if key, ok := skinparamKeys[property]; ok {
		if v, exists := r.skinparams[key]; exists {
			return v
		}
	}
	if v, exists := r.skinparams[property]; exists {
		return v
	}
	if v := r.themeColor(property); v != "" {
		return v
	}
	return r.fallbackColor(property)
}

// ResolveInt returns the integer value for a named property.
// Resolution order: skinparam → theme → fallback default.
func (r *Resolver) ResolveInt(property string, fallback int) int {
	if key, ok := skinparamKeys[property]; ok {
		if v, exists := r.skinparams[key]; exists {
			return atoiOr(v, fallback)
		}
	}
	if v, exists := r.skinparams[property]; exists {
		return atoiOr(v, fallback)
	}
	if v := intFieldByName(r.theme, property); v != 0 {
		return v
	}
	if v := intFieldByName(r.fallback, property); v != 0 {
		return v
	}
	return fallback
}

func (r *Resolver) themeColor(property string) string {
	return fieldByName(r.theme, property)
}

func (r *Resolver) fallbackColor(property string) string {
	return fieldByName(r.fallback, property)
}

// fieldByName returns the string value of a Theme field by name.
// Returns "" if the field doesn't exist or isn't a string.
func fieldByName(t *Theme, name string) string {
	switch name {
	case "BackgroundColor":
		return t.BackgroundColor
	case "FontName":
		return t.FontName
	case "FontColor":
		return t.FontColor
	case "ClassBackgroundColor":
		return t.ClassBackgroundColor
	case "ClassBorderColor":
		return t.ClassBorderColor
	case "ClassFontColor":
		return t.ClassFontColor
	case "ClassStereotypeFontColor":
		return t.ClassStereotypeFontColor
	case "InterfaceBackgroundColor":
		return t.InterfaceBackgroundColor
	case "InterfaceBorderColor":
		return t.InterfaceBorderColor
	case "InterfaceFontColor":
		return t.InterfaceFontColor
	case "EnumBackgroundColor":
		return t.EnumBackgroundColor
	case "EnumBorderColor":
		return t.EnumBorderColor
	case "EnumFontColor":
		return t.EnumFontColor
	case "ArrowColor":
		return t.ArrowColor
	case "NoteBackgroundColor":
		return t.NoteBackgroundColor
	case "NoteBorderColor":
		return t.NoteBorderColor
	case "NoteFontColor":
		return t.NoteFontColor
	case "ParticipantBackgroundColor":
		return t.ParticipantBackgroundColor
	case "ParticipantBorderColor":
		return t.ParticipantBorderColor
	case "ParticipantFontColor":
		return t.ParticipantFontColor
	case "SequenceLifeLineBorderColor":
		return t.SequenceLifeLineBorderColor
	case "PackageBackgroundColor":
		return t.PackageBackgroundColor
	case "PackageBorderColor":
		return t.PackageBorderColor
	case "PackageFontColor":
		return t.PackageFontColor
	case "AnnotationColor":
		return t.AnnotationColor
	default:
		return ""
	}
}

// intFieldByName returns the int value of a Theme field by name.
func intFieldByName(t *Theme, name string) int {
	switch name {
	case "FontSize":
		return t.FontSize
	case "ClassFontSize":
		return t.ClassFontSize
	case "ArrowFontSize":
		return t.ArrowFontSize
	case "Padding":
		return t.Padding
	case "ClassPadding":
		return t.ClassPadding
	case "NotePadding":
		return t.NotePadding
	case "BorderWidth":
		return t.BorderWidth
	case "ArrowThickness":
		return t.ArrowThickness
	default:
		return 0
	}
}

func atoiOr(s string, fallback int) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return fallback
		}
		n = n*10 + int(ch-'0')
	}
	if s == "" {
		return fallback
	}
	return n
}
