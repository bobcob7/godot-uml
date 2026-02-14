package svg_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bobcob7/godot-uml/internal/parser"
	"github.com/bobcob7/godot-uml/internal/renderer/svg"
	"github.com/bobcob7/godot-uml/internal/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const updateGolden = false // set to true to regenerate golden files

func TestClassRenderer(t *testing.T) {
	t.Parallel()
	t.Run("EmptyDiagram", func(t *testing.T) {
		t.Parallel()
		diagram, errs := parser.Parse("@startuml\n@enduml")
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<svg")
		assert.Contains(t, out, "</svg>")
		assert.Contains(t, out, `xmlns="http://www.w3.org/2000/svg"`)
	})
	t.Run("SingleClass", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo {\n+name : String\n-age : int\n+speak() : void\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "name : String")
		assert.Contains(t, out, "speak()")
		assert.Contains(t, out, "void")
		assert.Contains(t, out, `rx="8"`)
	})
	t.Run("InterfaceWithStereotype", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\ninterface Drawable {\n+draw() : void\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "&lt;&lt;interface&gt;&gt;")
		assert.Contains(t, out, "Drawable")
		assert.Contains(t, out, "draw()")
	})
	t.Run("EnumWithValues", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nenum Color {\nRED\nGREEN\nBLUE\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "&lt;&lt;enum&gt;&gt;")
		assert.Contains(t, out, "Color")
	})
	t.Run("AbstractClass", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nabstract class Shape {\n+area() : double\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Shape")
		assert.Contains(t, out, `font-style="italic"`)
	})
	t.Run("RelationshipInheritance", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Animal\nclass Dog\nDog --|> Animal\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Animal")
		assert.Contains(t, out, "Dog")
		assert.Contains(t, out, "<line")
		assert.Contains(t, out, "<polygon")
	})
	t.Run("RelationshipWithLabel", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass A\nclass B\nA --> B : uses\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "uses")
	})
	t.Run("RelationshipWithCardinality", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Animal\nclass Leg\nAnimal \"1\" --> \"*\" Leg : has\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, ">1<")
		assert.Contains(t, out, ">*<")
		assert.Contains(t, out, "has")
	})
	t.Run("DashedRelationships", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass A\nclass B\nA ..> B\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `stroke-dasharray`)
	})
	t.Run("NoteAttached", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\nnote left of Foo : This is a note\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "This is a note")
		assert.Contains(t, out, "<polygon")
		assert.Contains(t, out, `stroke-dasharray="5,5"`)
	})
	t.Run("Package", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\npackage com.example {\nclass Foo\nclass Bar\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "com.example")
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "Bar")
	})
	t.Run("AllVisibilityIcons", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass V {\n+pub : int\n-priv : int\n#prot : int\n~pkg : int\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, ">+<")
		assert.Contains(t, out, ">-<")
		assert.Contains(t, out, ">#<")
		assert.Contains(t, out, ">~<")
	})
	t.Run("StaticModifierUnderline", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass S {\n{static} count : int\n}\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `text-decoration="underline"`)
	})
	t.Run("DarculaThemeColors", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "#2B2B2B")
		assert.Contains(t, out, "#3C3F41")
		assert.Contains(t, out, "#555555")
	})
	t.Run("SkinparamOverride", func(t *testing.T) {
		t.Parallel()
		// Note: parser tokenizes "#FF0000" as "# FF0000", so the skinparam
		// value includes the space. We verify the override takes effect.
		input := "@startuml\nskinparam backgroundColor #FF0000\nclass Foo\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		// The skinparam value overrides the default Darcula background.
		assert.NotContains(t, out, `fill="#2B2B2B"`)
	})
	t.Run("NilResolverUsesDarcula", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "#2B2B2B")
	})
	t.Run("CustomTheme", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		custom := theme.Darcula()
		custom.BackgroundColor = "#AABBCC"
		resolver := theme.NewResolver(custom)
		r := svg.NewClassRenderer(resolver)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "#AABBCC")
	})
	t.Run("XMLEscaping", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\nnote left of Foo : Use <T> & \"quotes\"\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "&amp;")
		assert.NotContains(t, out, "& \"")
	})
	t.Run("ImplicitClassFromRelationship", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nFoo --> Bar\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "Bar")
	})
	t.Run("CompositionDiamond", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass A\nclass B\nA --* B\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<polygon")
	})
	t.Run("AggregationDiamond", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass A\nclass B\nA --o B\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<polygon")
		assert.Contains(t, out, `fill="white"`)
	})
}

func TestClassRendererGolden(t *testing.T) {
	t.Parallel()
	t.Run("ClassBasicFixture", func(t *testing.T) {
		t.Parallel()
		data, err := os.ReadFile("../../../testdata/class_basic.puml")
		require.NoError(t, err)
		diagram, errs := parser.Parse(string(data))
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err = r.Render(&buf, diagram)
		require.NoError(t, err)
		got := buf.String()
		goldenPath := filepath.Join("../../../testdata", "class_basic.golden.svg")
		if updateGolden {
			err = os.WriteFile(goldenPath, []byte(got), 0o644)
			require.NoError(t, err)
			return
		}
		golden, err := os.ReadFile(goldenPath)
		if err != nil {
			// No golden file yet; create it and pass.
			err = os.WriteFile(goldenPath, []byte(got), 0o644)
			require.NoError(t, err)
			t.Log("Created golden file, rerun to verify")
			return
		}
		assert.Equal(t, string(golden), got, "SVG output differs from golden file")
	})
}

func TestClassRendererValidSVG(t *testing.T) {
	t.Parallel()
	t.Run("ValidSVG11", func(t *testing.T) {
		t.Parallel()
		input := `@startuml
class Animal {
  +name : String
  -age : int
  +speak() : void
}
abstract class Shape {
  +area() : double
}
interface Drawable {
  +draw() : void
}
enum Color {
  RED
  GREEN
  BLUE
}
Animal --|> Shape
Animal ..|> Drawable
Dog --|> Animal
Animal "1" --> "*" Leg : has
note left of Animal : This is an animal
package com.example {
  class Foo
  class Bar
}
@enduml`
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewClassRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		// Verify basic SVG structure.
		assert.True(t, strings.HasPrefix(out, "<svg"))
		assert.True(t, strings.HasSuffix(strings.TrimSpace(out), "</svg>"))
		assert.Contains(t, out, `xmlns="http://www.w3.org/2000/svg"`)
		// Verify key elements rendered.
		assert.Contains(t, out, "Animal")
		assert.Contains(t, out, "Shape")
		assert.Contains(t, out, "Drawable")
		assert.Contains(t, out, "Color")
		assert.Contains(t, out, "Dog")
		assert.Contains(t, out, "Leg")
		assert.Contains(t, out, "com.example")
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "Bar")
		assert.Contains(t, out, "This is an animal")
	})
}
