package parser

import (
	"os"
	"testing"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParticipant(t *testing.T) {
	t.Parallel()
	t.Run("BasicParticipant", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Alice", p.Name)
		assert.Equal(t, ast.ParticipantDefault, p.Kind)
	})
	t.Run("Actor", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nactor Bob\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Bob", p.Name)
		assert.Equal(t, ast.ParticipantActor, p.Kind)
	})
	t.Run("Boundary", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nboundary Web\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Web", p.Name)
		assert.Equal(t, ast.ParticipantBoundary, p.Kind)
	})
	t.Run("Control", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ncontrol Router\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Router", p.Name)
		assert.Equal(t, ast.ParticipantControl, p.Kind)
	})
	t.Run("Entity", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nentity User\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "User", p.Name)
		assert.Equal(t, ast.ParticipantEntity, p.Kind)
	})
	t.Run("Database", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ndatabase DB\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "DB", p.Name)
		assert.Equal(t, ast.ParticipantDatabase, p.Kind)
	})
	t.Run("Collections", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ncollections Workers\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Workers", p.Name)
		assert.Equal(t, ast.ParticipantCollections, p.Kind)
	})
	t.Run("Queue", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nqueue Jobs\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Jobs", p.Name)
		assert.Equal(t, ast.ParticipantQueue, p.Kind)
	})
	t.Run("WithAlias", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant \"Long Name\" as LN\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		p, ok := diagram.Statements[0].(*ast.Participant)
		require.True(t, ok)
		assert.Equal(t, "Long Name", p.Name)
		assert.Equal(t, "LN", p.Alias)
	})
}

func TestParseMessage(t *testing.T) {
	t.Parallel()
	t.Run("SolidArrow", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Equal(t, "Alice", m.From)
		assert.Equal(t, "Bob", m.To)
		assert.Equal(t, "hello", m.Label)
		assert.Equal(t, "->", m.Arrow)
		assert.False(t, m.Dashed)
	})
	t.Run("DashedArrow", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Bob\nparticipant Alice\nBob --> Alice : response\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Equal(t, "Bob", m.From)
		assert.Equal(t, "Alice", m.To)
		assert.Equal(t, "response", m.Label)
		assert.True(t, m.Dashed)
	})
	t.Run("LeftArrow", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nactor Alice\nactor Bob\nAlice <- Bob : data\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Equal(t, "Alice", m.From)
		assert.Equal(t, "Bob", m.To)
		assert.Equal(t, "data", m.Label)
		assert.Equal(t, "<-", m.Arrow)
	})
	t.Run("DottedArrow", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nparticipant Bob\nAlice ..> Bob : async\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.True(t, m.Dashed)
	})
	t.Run("NoLabel", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Empty(t, m.Label)
	})
	t.Run("MultipleMessages", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\nBob --> Alice : world\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 4)
	})
	t.Run("ActivationShorthandPlusPlus", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob ++ : activate\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Equal(t, "activate", m.Label)
	})
	t.Run("ActivationShorthandMinusMinus", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Bob\nparticipant Alice\nBob -> Alice -- : deactivate\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		m, ok := diagram.Statements[2].(*ast.Message)
		require.True(t, ok)
		assert.Equal(t, "deactivate", m.Label)
	})
}

func TestIsSequenceArrow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arrow string
		want  bool
	}{
		{"RightSolid", "->", true},
		{"LeftSolid", "<-", true},
		{"Bidirectional", "<->", true},
		{"DashedRight", "-->", false},
		{"DashedLeft", "<--", false},
		{"Inheritance", "--|>", false},
		{"Realization", "..|>", false},
		{"Composition", "--*", false},
		{"Aggregation", "--o", false},
		{"Dependency", "..>", false},
		{"PlainDouble", "--", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isSequenceArrow(tt.arrow))
		})
	}
}

func TestImplicitSequenceMode(t *testing.T) {
	t.Parallel()
	t.Run("SolidArrowWithoutParticipant", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nBob -> Alice : hello\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		m, ok := diagram.Statements[0].(*ast.Message)
		require.True(t, ok, "expected *ast.Message, got %T", diagram.Statements[0])
		assert.Equal(t, "Bob", m.From)
		assert.Equal(t, "Alice", m.To)
		assert.Equal(t, "hello", m.Label)
		assert.False(t, m.Dashed)
	})
	t.Run("DashedArrowAfterSolid", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nAlice -> Bob : hi\nBob --> Alice : bye\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		m1, ok := diagram.Statements[0].(*ast.Message)
		require.True(t, ok, "stmt 0: expected *ast.Message, got %T", diagram.Statements[0])
		assert.Equal(t, "Alice", m1.From)
		assert.False(t, m1.Dashed)
		m2, ok := diagram.Statements[1].(*ast.Message)
		require.True(t, ok, "stmt 1: expected *ast.Message, got %T", diagram.Statements[1])
		assert.Equal(t, "Bob", m2.From)
		assert.True(t, m2.Dashed)
	})
	t.Run("FullDiagramWithNotes", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nBob -> Alice: hello\nAlice -> Tom: Bob says hello\nnote over Alice\nShoots Tom\nend note\nAlice --> Bob: done\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 4)
		_, ok := diagram.Statements[0].(*ast.Message)
		assert.True(t, ok, "stmt 0: expected *ast.Message, got %T", diagram.Statements[0])
		_, ok = diagram.Statements[1].(*ast.Message)
		assert.True(t, ok, "stmt 1: expected *ast.Message, got %T", diagram.Statements[1])
		n, ok := diagram.Statements[2].(*ast.Note)
		require.True(t, ok, "stmt 2: expected *ast.Note, got %T", diagram.Statements[2])
		assert.Equal(t, "Alice", n.Target)
		assert.Equal(t, "Shoots Tom", n.Text)
		m, ok := diagram.Statements[3].(*ast.Message)
		require.True(t, ok, "stmt 3: expected *ast.Message, got %T", diagram.Statements[3])
		assert.Equal(t, "Alice", m.From)
		assert.Equal(t, "Bob", m.To)
		assert.True(t, m.Dashed)
	})
	t.Run("ClassRelationshipUnchanged", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nFoo --> Bar : uses\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		_, ok := diagram.Statements[0].(*ast.Relationship)
		assert.True(t, ok, "expected *ast.Relationship for --> arrow, got %T", diagram.Statements[0])
	})
}

func TestParseActivate(t *testing.T) {
	t.Parallel()
	t.Run("Activate", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nactivate Bob\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		a, ok := diagram.Statements[0].(*ast.Activate)
		require.True(t, ok)
		assert.Equal(t, "Bob", a.Target)
		assert.False(t, a.Deactivate)
	})
	t.Run("Deactivate", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ndeactivate Bob\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		a, ok := diagram.Statements[0].(*ast.Activate)
		require.True(t, ok)
		assert.Equal(t, "Bob", a.Target)
		assert.True(t, a.Deactivate)
	})
}

func TestParseReturn(t *testing.T) {
	t.Parallel()
	t.Run("WithLabel", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nreturn success\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		r, ok := diagram.Statements[0].(*ast.Return)
		require.True(t, ok)
		assert.Equal(t, "success", r.Label)
	})
	t.Run("WithoutLabel", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nreturn\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		r, ok := diagram.Statements[0].(*ast.Return)
		require.True(t, ok)
		assert.Empty(t, r.Label)
	})
}

func TestParseFragment(t *testing.T) {
	t.Parallel()
	t.Run("AltElse", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nalt success\nAlice -> Bob : ok\nelse failure\nAlice -> Bob : error\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentAlt, f.Kind)
		assert.Equal(t, "success", f.Condition)
		require.Len(t, f.Statements, 1)
		require.Len(t, f.ElseParts, 1)
		assert.Equal(t, "failure", f.ElseParts[0].Condition)
		require.Len(t, f.ElseParts[0].Statements, 1)
	})
	t.Run("Loop", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nloop 10 times\nAlice -> Bob : ping\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentLoop, f.Kind)
		assert.Equal(t, "10 times", f.Condition)
		require.Len(t, f.Statements, 1)
	})
	t.Run("Par", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nparticipant Charlie\npar\nAlice -> Bob : task1\nelse\nAlice -> Charlie : task2\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 4)
		f, ok := diagram.Statements[3].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentPar, f.Kind)
		require.Len(t, f.Statements, 1)
		require.Len(t, f.ElseParts, 1)
	})
	t.Run("Group", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\ngroup My Group\nAlice -> Bob : msg\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentGroup, f.Kind)
		assert.Equal(t, "My Group", f.Condition)
	})
	t.Run("Break", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nbreak emergency\nAlice -> Bob : stop\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentBreak, f.Kind)
	})
	t.Run("Ref", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nref over Alice\nAlice -> Bob : see other\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentRef, f.Kind)
		assert.Equal(t, "over Alice", f.Condition)
	})
	t.Run("NestedFragments", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nalt outer\nloop 3 times\nAlice -> Bob : msg\nend\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		outer, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentAlt, outer.Kind)
		require.Len(t, outer.Statements, 1)
		inner, ok := outer.Statements[0].(*ast.Fragment)
		require.True(t, ok)
		assert.Equal(t, ast.FragmentLoop, inner.Kind)
	})
	t.Run("MultipleElseParts", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nparticipant Bob\nalt case1\nAlice -> Bob : a\nelse case2\nAlice -> Bob : b\nelse case3\nAlice -> Bob : c\nend\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		f, ok := diagram.Statements[2].(*ast.Fragment)
		require.True(t, ok)
		require.Len(t, f.ElseParts, 2)
		assert.Equal(t, "case2", f.ElseParts[0].Condition)
		assert.Equal(t, "case3", f.ElseParts[1].Condition)
	})
}

func TestParseSequenceNote(t *testing.T) {
	t.Parallel()
	t.Run("LeftOfParticipant", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nnote left of Alice : Hello\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteLeft, n.Placement)
		assert.Equal(t, "Alice", n.Target)
		assert.Equal(t, "Hello", n.Text)
	})
	t.Run("RightOfParticipant", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Bob\nnote right of Bob : World\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteRight, n.Placement)
		assert.Equal(t, "Bob", n.Target)
		assert.Equal(t, "World", n.Text)
	})
	t.Run("OverParticipant", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nnote over Alice : Note text\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteOver, n.Placement)
		assert.Equal(t, "Alice", n.Target)
	})
	t.Run("OverMultipleParticipants", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nparticipant Alice\nparticipant Bob\nnote over Alice,Bob : Shared\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 3)
		n, ok := diagram.Statements[2].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, "Alice,Bob", n.Target)
	})
	t.Run("MultiLine", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nparticipant Alice\nnote left of Alice\nLine 1\nLine 2\nend note\n@enduml"
		diagram, errs := Parse(input)
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 2)
		n, ok := diagram.Statements[1].(*ast.Note)
		require.True(t, ok)
		assert.Contains(t, n.Text, "Line 1")
		assert.Contains(t, n.Text, "Line 2")
	})
}

func TestParseAutonumber(t *testing.T) {
	t.Parallel()
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nautonumber\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		a, ok := diagram.Statements[0].(*ast.Autonumber)
		require.True(t, ok)
		assert.Empty(t, a.Start)
	})
	t.Run("WithStart", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nautonumber 10\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		a, ok := diagram.Statements[0].(*ast.Autonumber)
		require.True(t, ok)
		assert.Equal(t, "10", a.Start)
	})
}

func TestParseDivider(t *testing.T) {
	t.Parallel()
	t.Run("WithText", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\n== Initialization ==\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		d, ok := diagram.Statements[0].(*ast.Divider)
		require.True(t, ok)
		assert.Equal(t, "Initialization", d.Text)
	})
}

func TestParseDelay(t *testing.T) {
	t.Parallel()
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\n...\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		d, ok := diagram.Statements[0].(*ast.Delay)
		require.True(t, ok)
		assert.Empty(t, d.Text)
	})
	t.Run("WithText", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\n... 5 minutes later ...\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		d, ok := diagram.Statements[0].(*ast.Delay)
		require.True(t, ok)
		assert.Equal(t, "5 minutes later", d.Text)
	})
}

func TestParseSequenceFixture(t *testing.T) {
	t.Parallel()
	t.Run("SequenceBasic", func(t *testing.T) {
		t.Parallel()
		data, err := os.ReadFile("../../testdata/sequence_basic.puml")
		require.NoError(t, err)
		diagram, errs := Parse(string(data))
		require.Empty(t, errs, "fixture should parse without errors: %v", errs)
		assert.NotEmpty(t, diagram.Statements)
		var participants, messages, fragments, notes int
		for _, stmt := range diagram.Statements {
			switch stmt.(type) {
			case *ast.Participant:
				participants++
			case *ast.Message:
				messages++
			case *ast.Fragment:
				fragments++
			case *ast.Note:
				notes++
			}
		}
		assert.Equal(t, 3, participants, "should have 3 participant declarations")
		assert.GreaterOrEqual(t, messages, 5, "should have at least 5 top-level messages")
		assert.Equal(t, 2, fragments, "should have 2 fragments (alt, loop)")
		assert.Equal(t, 3, notes, "should have 3 notes")
	})
}
