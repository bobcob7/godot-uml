package gouml_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/bobcob7/go-uml/internal/theme"
	"github.com/bobcob7/go-uml/pkg/gouml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	t.Parallel()
	t.Run("ClassDiagram", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo {\n+name : String\n}\n@enduml")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<svg")
		assert.Contains(t, out, "</svg>")
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "name")
	})
	t.Run("SequenceDiagram", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\n@enduml")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<svg")
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Bob")
		assert.Contains(t, out, "hello")
	})
	t.Run("EmptyDiagram", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\n@enduml")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "<svg")
	})
	t.Run("WithSkinparam", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo\n@enduml")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf,
			gouml.WithSkinparam("backgroundColor", "#FF0000"),
		)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "#FF0000")
	})
	t.Run("WithCustomTheme", func(t *testing.T) {
		t.Parallel()
		custom := theme.Darcula()
		custom.BackgroundColor = "#AABBCC"
		input := strings.NewReader("@startuml\nclass Foo\n@enduml")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf, gouml.WithTheme(custom))
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "#AABBCC")
	})
	t.Run("InvalidInput", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("not a diagram")
		var buf bytes.Buffer
		err := gouml.Render(input, &buf)
		require.Error(t, err)
	})
	t.Run("ReadError", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		err := gouml.Render(&errReader{}, &buf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reading input")
	})
}

func TestParse(t *testing.T) {
	t.Parallel()
	t.Run("ValidDiagram", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo\n@enduml")
		diagram, errs := gouml.Parse(input)
		require.Empty(t, errs)
		require.NotNil(t, diagram)
	})
	t.Run("InvalidDiagram", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("not a diagram")
		diagram, errs := gouml.Parse(input)
		require.NotEmpty(t, errs)
		require.NotNil(t, diagram)
	})
	t.Run("ErrorHasPosition", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("not a diagram")
		_, errs := gouml.Parse(input)
		require.NotEmpty(t, errs)
		assert.Greater(t, errs[0].Line, 0)
		assert.NotEmpty(t, errs[0].Message)
		assert.Contains(t, errs[0].Error(), errs[0].Message)
	})
	t.Run("RenderAfterParse", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo\n@enduml")
		diagram, errs := gouml.Parse(input)
		require.Empty(t, errs)
		var buf bytes.Buffer
		err := gouml.RenderDiagram(&buf, diagram)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Foo")
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()
	t.Run("ValidInput", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo\n@enduml")
		errs := gouml.Validate(input)
		assert.Empty(t, errs)
	})
	t.Run("InvalidInput", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("not a diagram")
		errs := gouml.Validate(input)
		assert.NotEmpty(t, errs)
	})
}

func TestRenderDiagram(t *testing.T) {
	t.Parallel()
	t.Run("SequenceDetection", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nparticipant Alice\nAlice -> Bob : hi\n@enduml")
		diagram, errs := gouml.Parse(input)
		require.Empty(t, errs)
		var buf bytes.Buffer
		err := gouml.RenderDiagram(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		// Sequence diagrams have lifeline dashes
		assert.Contains(t, out, `stroke-dasharray="5,5"`)
	})
	t.Run("ImplicitSequenceDetection", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nBob -> Alice: hello\nAlice -> Tom: Bob says hello\nnote over Alice\nShoots Tom\nend note\nAlice --> Bob: done\n@enduml")
		diagram, errs := gouml.Parse(input)
		require.Empty(t, errs)
		var buf bytes.Buffer
		err := gouml.RenderDiagram(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `stroke-dasharray="5,5"`, "should have lifeline dashes")
		assert.Contains(t, out, "Bob")
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Tom")
		assert.Contains(t, out, "hello")
		assert.Contains(t, out, "Shoots Tom")
		assert.Contains(t, out, "done")
	})
	t.Run("ClassDetection", func(t *testing.T) {
		t.Parallel()
		input := strings.NewReader("@startuml\nclass Foo {\n+name : String\n}\n@enduml")
		diagram, errs := gouml.Parse(input)
		require.Empty(t, errs)
		var buf bytes.Buffer
		err := gouml.RenderDiagram(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Foo")
		assert.Contains(t, out, "name")
	})
}

// errReader is an io.Reader that always returns an error.
type errReader struct{}

func (e *errReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
