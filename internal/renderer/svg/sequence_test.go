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

const updateSeqGolden = false

func TestSequenceRenderer(t *testing.T) {
	t.Parallel()
	t.Run("EmptyDiagram", func(t *testing.T) {
		t.Parallel()
		diagram, errs := parser.Parse("@startuml\n@enduml")
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "<svg")
		assert.Contains(t, out, "</svg>")
		assert.Contains(t, out, `xmlns="http://www.w3.org/2000/svg"`)
	})
	t.Run("ParticipantBoxes", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Bob")
		assert.Contains(t, out, "<rect")
	})
	t.Run("ActorParticipant", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nactor Bob\nparticipant Alice\nAlice -> Bob : hello\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Bob")
		assert.Contains(t, out, "<circle")
	})
	t.Run("SyncMessage", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : request\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "request")
		assert.Contains(t, out, "<line")
		assert.Contains(t, out, "<polygon")
	})
	t.Run("AsyncDashedMessage", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice --> Bob : response\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "response")
		assert.Contains(t, out, `stroke-dasharray`)
	})
	t.Run("LifelineDashes", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hi\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, `stroke-dasharray="5,5"`)
	})
	t.Run("ActivationBars", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nactivate Bob\nAlice -> Bob : work\ndeactivate Bob\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		// Activation bar is a narrow rect on the lifeline
		rectCount := strings.Count(out, "<rect")
		// At least: bg rect + 2 top participant boxes + 2 bottom participant boxes + 1 activation
		assert.GreaterOrEqual(t, rectCount, 6)
	})
	t.Run("NoteLeft", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote left of Alice : Client side\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Client side")
		assert.Contains(t, out, "<polygon")
	})
	t.Run("NoteRight", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote right of Alice : Server side\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Server side")
	})
	t.Run("NoteOver", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote over Alice : Centered note\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Centered note")
	})
	t.Run("AltElseFragment", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nalt success\nAlice -> Bob : ok\nelse failure\nAlice -> Bob : retry\nend\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "alt")
		assert.Contains(t, out, "success")
		assert.Contains(t, out, "else")
		assert.Contains(t, out, "failure")
	})
	t.Run("LoopFragment", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nloop 3 times\nAlice -> Bob : poll\nend\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "loop")
		assert.Contains(t, out, "3 times")
	})
	t.Run("Divider", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : first\n== Phase 2 ==\nAlice -> Bob : second\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Phase 2")
	})
	t.Run("Delay", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : start\n... 5 minutes later ...\nAlice -> Bob : resume\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "5 minutes later")
		assert.Contains(t, out, `font-style="italic"`)
	})
	t.Run("Autonumber", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nautonumber\nAlice -> Bob : first\nAlice -> Bob : second\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "1.")
		assert.Contains(t, out, "2.")
	})
	t.Run("DarculaThemeColors", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nAlice -> Alice : self\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "#2B2B2B")
		assert.Contains(t, out, "#3C3F41")
	})
	t.Run("CustomTheme", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		custom := theme.Darcula()
		custom.BackgroundColor = "#AABBCC"
		resolver := theme.NewResolver(custom)
		r := svg.NewSequenceRenderer(resolver)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "#AABBCC")
	})
	t.Run("ImplicitParticipants", func(t *testing.T) {
		t.Parallel()
		// participant keyword triggers seqMode; then Bob is implicit from message
		input := "@startuml\nparticipant Alice\nAlice -> Bob : hello\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Bob")
	})
	t.Run("BottomParticipantBoxes", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : msg\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		// Should have participant boxes at top AND bottom (4 rects for 2 participants + 1 bg)
		rectCount := strings.Count(out, "<rect")
		assert.GreaterOrEqual(t, rectCount, 5)
	})
	t.Run("DatabaseParticipant", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\ndatabase DB\nparticipant Alice\nAlice -> DB : query\n@enduml"
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "DB")
		assert.Contains(t, out, "query")
	})
}

func TestSequenceRendererGolden(t *testing.T) {
	t.Parallel()
	t.Run("SequenceBasicFixture", func(t *testing.T) {
		t.Parallel()
		data, err := os.ReadFile("../../../testdata/sequence_basic.puml")
		require.NoError(t, err)
		diagram, errs := parser.Parse(string(data))
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err = r.Render(&buf, diagram)
		require.NoError(t, err)
		got := buf.String()
		goldenPath := filepath.Join("../../../testdata", "sequence_basic.golden.svg")
		if updateSeqGolden {
			err = os.WriteFile(goldenPath, []byte(got), 0o644)
			require.NoError(t, err)
			return
		}
		golden, err := os.ReadFile(goldenPath)
		if err != nil {
			err = os.WriteFile(goldenPath, []byte(got), 0o644)
			require.NoError(t, err)
			t.Log("Created golden file, rerun to verify")
			return
		}
		assert.Equal(t, string(golden), got, "SVG output differs from golden file")
	})
}

func TestSequenceRendererValidSVG(t *testing.T) {
	t.Parallel()
	t.Run("ValidSVG11", func(t *testing.T) {
		t.Parallel()
		input := `@startuml
participant Alice
actor Bob
database DB

Alice -> Bob : authenticate
Bob -> DB : query
DB --> Bob : result
Bob --> Alice : response

note left of Alice : Client
note right of Bob : Server
note over DB : Storage

alt success
  Alice -> Bob : confirmed
else failure
  Alice -> Bob : retry
end

loop 3 times
  Bob -> DB : poll
end

== Phase 2 ==

autonumber

... 5 minutes later ...

activate Bob
Alice -> Bob : resume
deactivate Bob
@enduml`
		diagram, errs := parser.Parse(input)
		require.Empty(t, errs)
		r := svg.NewSequenceRenderer(nil)
		var buf bytes.Buffer
		err := r.Render(&buf, diagram)
		require.NoError(t, err)
		out := buf.String()
		assert.True(t, strings.HasPrefix(out, "<svg"))
		assert.True(t, strings.HasSuffix(strings.TrimSpace(out), "</svg>"))
		assert.Contains(t, out, `xmlns="http://www.w3.org/2000/svg"`)
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "Bob")
		assert.Contains(t, out, "DB")
		assert.Contains(t, out, "authenticate")
		assert.Contains(t, out, "Phase 2")
		assert.Contains(t, out, "5 minutes later")
		assert.Contains(t, out, "alt")
		assert.Contains(t, out, "loop")
	})
}
