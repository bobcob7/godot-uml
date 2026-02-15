package parser

import (
	"os"
	"testing"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseClassDef(t *testing.T) {
	t.Parallel()
	t.Run("BasicClass", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass Foo {\n}\n@enduml")
		require.Empty(t, errs)
		require.Len(t, diagram.Statements, 1)
		cd, ok := diagram.Statements[0].(*ast.ClassDef)
		require.True(t, ok)
		assert.Equal(t, "Foo", cd.Name)
		assert.False(t, cd.Abstract)
	})
	t.Run("ClassWithFields", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass Foo {\n+name : String\n-age : int\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		require.Len(t, cd.Members, 2)
		f1, ok := cd.Members[0].(*ast.Field)
		require.True(t, ok)
		assert.Equal(t, "name", f1.Name)
		assert.Equal(t, "String", f1.Type)
		assert.Equal(t, ast.VisibilityPublic, f1.Visibility)
		f2 := cd.Members[1].(*ast.Field)
		assert.Equal(t, "age", f2.Name)
		assert.Equal(t, ast.VisibilityPrivate, f2.Visibility)
	})
	t.Run("ClassWithMethods", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass Foo {\n+speak() : void\n-calc(x : int) : int\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		require.Len(t, cd.Members, 2)
		m1, ok := cd.Members[0].(*ast.Method)
		require.True(t, ok)
		assert.Equal(t, "speak", m1.Name)
		assert.Equal(t, "void", m1.ReturnType)
		assert.Equal(t, ast.VisibilityPublic, m1.Visibility)
		m2 := cd.Members[1].(*ast.Method)
		assert.Equal(t, "calc", m2.Name)
		assert.Equal(t, "x : int", m2.Params)
		assert.Equal(t, "int", m2.ReturnType)
	})
	t.Run("AllVisibilityModifiers", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass V {\n+pub : int\n-priv : int\n#prot : int\n~pkg : int\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		require.Len(t, cd.Members, 4)
		assert.Equal(t, ast.VisibilityPublic, cd.Members[0].(*ast.Field).Visibility)
		assert.Equal(t, ast.VisibilityPrivate, cd.Members[1].(*ast.Field).Visibility)
		assert.Equal(t, ast.VisibilityProtected, cd.Members[2].(*ast.Field).Visibility)
		assert.Equal(t, ast.VisibilityPackage, cd.Members[3].(*ast.Field).Visibility)
	})
	t.Run("StaticModifier", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass S {\n{static} count : int\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		require.Len(t, cd.Members, 1)
		f := cd.Members[0].(*ast.Field)
		assert.Equal(t, ast.ModifierStatic, f.Modifier)
	})
	t.Run("AbstractClass", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nabstract class Shape {\n+area() : double\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		assert.Equal(t, "Shape", cd.Name)
		assert.True(t, cd.Abstract)
	})
	t.Run("AbstractAlone", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nabstract Shape\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		assert.Equal(t, "Shape", cd.Name)
		assert.True(t, cd.Abstract)
	})
	t.Run("ClassWithStereotype", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass Foo <<service>> {\n}\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		assert.Equal(t, "service", cd.Stereotype)
	})
	t.Run("ClassNoBody", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass Foo\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		assert.Equal(t, "Foo", cd.Name)
		assert.Empty(t, cd.Members)
	})
	t.Run("ClassWithAlias", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nclass \"Long Name\" as LN\n@enduml")
		require.Empty(t, errs)
		cd := diagram.Statements[0].(*ast.ClassDef)
		assert.Equal(t, "Long Name", cd.Name)
		assert.Equal(t, "LN", cd.Alias)
	})
}

func TestParseInterfaceDef(t *testing.T) {
	t.Parallel()
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ninterface Drawable {\n+draw() : void\n}\n@enduml")
		require.Empty(t, errs)
		idef, ok := diagram.Statements[0].(*ast.InterfaceDef)
		require.True(t, ok)
		assert.Equal(t, "Drawable", idef.Name)
		require.Len(t, idef.Members, 1)
	})
	t.Run("NoBody", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\ninterface Runnable\n@enduml")
		require.Empty(t, errs)
		idef := diagram.Statements[0].(*ast.InterfaceDef)
		assert.Equal(t, "Runnable", idef.Name)
	})
}

func TestParseEnumDef(t *testing.T) {
	t.Parallel()
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nenum Color {\nRED\nGREEN\nBLUE\n}\n@enduml")
		require.Empty(t, errs)
		edef, ok := diagram.Statements[0].(*ast.EnumDef)
		require.True(t, ok)
		assert.Equal(t, "Color", edef.Name)
		require.Len(t, edef.Members, 3)
	})
}

func TestParseRelationship(t *testing.T) {
	t.Parallel()
	t.Run("Inheritance", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nDog --|> Animal\n@enduml")
		require.Empty(t, errs)
		rel, ok := diagram.Statements[0].(*ast.Relationship)
		require.True(t, ok)
		assert.Equal(t, "Dog", rel.Left)
		assert.Equal(t, "Animal", rel.Right)
		assert.Equal(t, ast.RelInheritance, rel.Type)
		assert.Equal(t, "--|>", rel.Arrow)
	})
	t.Run("Realization", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nFoo ..|> Bar\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelRealization, rel.Type)
	})
	t.Run("Composition", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA --* B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelComposition, rel.Type)
	})
	t.Run("Aggregation", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA --o B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelAggregation, rel.Type)
	})
	t.Run("Association", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA --> B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelAssociation, rel.Type)
		assert.Equal(t, ast.ArrowRight, rel.Direction)
	})
	t.Run("LeftArrow", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA <-- B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.ArrowLeft, rel.Direction)
	})
	t.Run("WithLabel", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nDog --|> Animal : extends\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, "extends", rel.Label)
	})
	t.Run("WithCardinality", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nAnimal \"1\" --> \"*\" Leg : has\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, "1", rel.LeftCard)
		assert.Equal(t, "*", rel.RightCard)
		assert.Equal(t, "has", rel.Label)
	})
	t.Run("Dependency", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA ..> B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelDependency, rel.Type)
	})
	t.Run("LeftInheritance", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nA <|-- B\n@enduml")
		require.Empty(t, errs)
		rel := diagram.Statements[0].(*ast.Relationship)
		assert.Equal(t, ast.RelInheritance, rel.Type)
		assert.Equal(t, ast.ArrowLeft, rel.Direction)
	})
}

func TestParsePackage(t *testing.T) {
	t.Parallel()
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\npackage com.example {\nclass Foo\n}\n@enduml")
		require.Empty(t, errs)
		pkg, ok := diagram.Statements[0].(*ast.Package)
		require.True(t, ok)
		assert.Equal(t, "com.example", pkg.Name)
		assert.False(t, pkg.IsNamespace)
		require.Len(t, pkg.Statements, 1)
	})
	t.Run("Namespace", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nnamespace com.example {\nclass Bar\n}\n@enduml")
		require.Empty(t, errs)
		pkg := diagram.Statements[0].(*ast.Package)
		assert.True(t, pkg.IsNamespace)
	})
}

func TestParseNote(t *testing.T) {
	t.Parallel()
	t.Run("LeftOf", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nnote left of Foo : hello\n@enduml")
		require.Empty(t, errs)
		n, ok := diagram.Statements[0].(*ast.Note)
		require.True(t, ok)
		assert.Equal(t, ast.NoteLeft, n.Placement)
		assert.Equal(t, "Foo", n.Target)
		assert.Equal(t, "hello", n.Text)
	})
	t.Run("RightOf", func(t *testing.T) {
		t.Parallel()
		diagram, errs := Parse("@startuml\nnote right of Bar : world\n@enduml")
		require.Empty(t, errs)
		n := diagram.Statements[0].(*ast.Note)
		assert.Equal(t, ast.NoteRight, n.Placement)
		assert.Equal(t, "Bar", n.Target)
	})
}

func TestParseFixture(t *testing.T) {
	t.Parallel()
	t.Run("ClassBasic", func(t *testing.T) {
		t.Parallel()
		data, err := os.ReadFile("../../testdata/class_basic.puml")
		require.NoError(t, err)
		diagram, errs := Parse(string(data))
		require.Empty(t, errs, "fixture should parse without errors: %v", errs)
		assert.NotEmpty(t, diagram.Statements, "should have parsed statements")
		var classes, interfaces, enums, relationships, packages int
		for _, stmt := range diagram.Statements {
			switch stmt.(type) {
			case *ast.ClassDef:
				classes++
			case *ast.InterfaceDef:
				interfaces++
			case *ast.EnumDef:
				enums++
			case *ast.Relationship:
				relationships++
			case *ast.Package:
				packages++
			}
		}
		assert.Equal(t, 2, classes, "Animal + Shape")
		assert.Equal(t, 1, interfaces, "Drawable")
		assert.Equal(t, 1, enums, "Color")
		assert.GreaterOrEqual(t, relationships, 4, "multiple relationships")
		assert.Equal(t, 1, packages, "com.example")
	})
}
