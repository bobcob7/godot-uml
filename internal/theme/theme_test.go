package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDarcula_PaletteValues(t *testing.T) {
	t.Parallel()
	d := Darcula()
	require.NotNil(t, d)
	assert.Equal(t, "#2B2B2B", d.BackgroundColor)
	assert.Equal(t, "#A9B7C6", d.FontColor)
	assert.Equal(t, "#3C3F41", d.ClassBackgroundColor)
	assert.Equal(t, "#555555", d.ClassBorderColor)
	assert.Equal(t, "#A9B7C6", d.ArrowColor)
	assert.Equal(t, "#4E5254", d.NoteBackgroundColor)
	assert.Equal(t, "#CC7832", d.ClassStereotypeFontColor)
	assert.Equal(t, "#6A8759", d.AnnotationColor)
	assert.Equal(t, "#6897BB", d.InterfaceFontColor)
}

func TestDarcula_FontDefaults(t *testing.T) {
	t.Parallel()
	d := Darcula()
	assert.Equal(t, "DejaVu Sans", d.FontName)
	assert.Equal(t, 13, d.FontSize)
	assert.Equal(t, 13, d.ClassFontSize)
}

func TestDarcula_SpacingDefaults(t *testing.T) {
	t.Parallel()
	d := Darcula()
	assert.Equal(t, 10, d.Padding)
	assert.Equal(t, 8, d.ClassPadding)
	assert.Equal(t, 8, d.NotePadding)
	assert.Equal(t, 1, d.BorderWidth)
	assert.Equal(t, 1, d.ArrowThickness)
}

func TestNewResolver_NilThemeUsesDarcula(t *testing.T) {
	t.Parallel()
	r := NewResolver(nil)
	assert.Equal(t, "#2B2B2B", r.ResolveColor("BackgroundColor"))
}

func TestNewResolver_CustomTheme(t *testing.T) {
	t.Parallel()
	custom := &Theme{
		BackgroundColor: "#FF0000",
	}
	r := NewResolver(custom)
	assert.Equal(t, "#FF0000", r.ResolveColor("BackgroundColor"))
}

func TestResolver_SkinparamOverridesTheme(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	r.SetSkinparam("backgroundColor", "#000000")
	assert.Equal(t, "#000000", r.ResolveColor("BackgroundColor"))
}

func TestResolver_SkinparamDirectKey(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	r.SetSkinparam("ClassBackgroundColor", "#112233")
	assert.Equal(t, "#112233", r.ResolveColor("ClassBackgroundColor"))
}

func TestResolver_ThemeOverridesFallback(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	// Darcula has #2B2B2B, fallback has #FFFFFF
	assert.Equal(t, "#2B2B2B", r.ResolveColor("BackgroundColor"))
}

func TestResolver_FallbackUsedWhenThemeEmpty(t *testing.T) {
	t.Parallel()
	emptyTheme := &Theme{}
	r := NewResolver(emptyTheme)
	assert.Equal(t, "#FFFFFF", r.ResolveColor("BackgroundColor"))
}

func TestResolver_UnknownProperty(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	assert.Equal(t, "", r.ResolveColor("nonexistentProperty"))
}

func TestResolver_MultipleSkinparams(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	r.SetSkinparam("backgroundColor", "#111111")
	r.SetSkinparam("classBorderColor", "#222222")
	assert.Equal(t, "#111111", r.ResolveColor("BackgroundColor"))
	assert.Equal(t, "#222222", r.ResolveColor("ClassBorderColor"))
	// Unset properties still resolve from theme.
	assert.Equal(t, "#A9B7C6", r.ResolveColor("FontColor"))
}

func TestResolver_SkinparamKeyMapping(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		spKey    string
		property string
		value    string
	}{
		{"background", "backgroundColor", "BackgroundColor", "#AABBCC"},
		{"class fill", "classBackgroundColor", "ClassBackgroundColor", "#112233"},
		{"arrow", "arrowColor", "ArrowColor", "#445566"},
		{"note fill", "noteBackgroundColor", "NoteBackgroundColor", "#778899"},
		{"font name", "defaultFontName", "FontName", "Courier"},
		{"annotation", "annotationColor", "AnnotationColor", "#ABCDEF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewResolver(Darcula())
			r.SetSkinparam(tt.spKey, tt.value)
			assert.Equal(t, tt.value, r.ResolveColor(tt.property))
		})
	}
}

func TestResolver_PrecedenceOrder(t *testing.T) {
	t.Parallel()
	// Verify: skinparam > theme > fallback
	r := NewResolver(&Theme{ClassBackgroundColor: "#THEME"})
	// Without skinparam, theme value is used.
	assert.Equal(t, "#THEME", r.ResolveColor("ClassBackgroundColor"))
	// With skinparam, skinparam wins.
	r.SetSkinparam("classBackgroundColor", "#SKINPARAM")
	assert.Equal(t, "#SKINPARAM", r.ResolveColor("ClassBackgroundColor"))
}

func TestResolver_AllDarculaColors(t *testing.T) {
	t.Parallel()
	r := NewResolver(Darcula())
	tests := []struct {
		property string
		want     string
	}{
		{"BackgroundColor", "#2B2B2B"},
		{"FontColor", "#A9B7C6"},
		{"ClassBackgroundColor", "#3C3F41"},
		{"ClassBorderColor", "#555555"},
		{"ClassFontColor", "#A9B7C6"},
		{"ClassStereotypeFontColor", "#CC7832"},
		{"InterfaceBackgroundColor", "#3C3F41"},
		{"InterfaceBorderColor", "#555555"},
		{"InterfaceFontColor", "#6897BB"},
		{"EnumBackgroundColor", "#3C3F41"},
		{"EnumBorderColor", "#555555"},
		{"EnumFontColor", "#A9B7C6"},
		{"ArrowColor", "#A9B7C6"},
		{"NoteBackgroundColor", "#4E5254"},
		{"NoteBorderColor", "#555555"},
		{"NoteFontColor", "#A9B7C6"},
		{"ParticipantBackgroundColor", "#3C3F41"},
		{"ParticipantBorderColor", "#555555"},
		{"ParticipantFontColor", "#A9B7C6"},
		{"SequenceLifeLineBorderColor", "#555555"},
		{"PackageBackgroundColor", "#2B2B2B"},
		{"PackageBorderColor", "#555555"},
		{"PackageFontColor", "#A9B7C6"},
		{"AnnotationColor", "#6A8759"},
	}
	for _, tt := range tests {
		t.Run(tt.property, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, r.ResolveColor(tt.property))
		})
	}
}

func TestHardcodedFallback_HasValues(t *testing.T) {
	t.Parallel()
	f := hardcodedFallback()
	require.NotNil(t, f)
	assert.Equal(t, "#FFFFFF", f.BackgroundColor)
	assert.Equal(t, "#000000", f.FontColor)
	assert.Equal(t, "sans-serif", f.FontName)
	assert.Equal(t, 12, f.FontSize)
}
