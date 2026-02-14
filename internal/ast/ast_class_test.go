package ast_test

import (
	"testing"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func TestClassDefStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 2, Column: 1}
		cd := &ast.ClassDef{Pos: pos, Name: "Foo"}
		var s ast.Statement = cd
		assert.Equal(t, pos, s.Position())
	})
}

func TestInterfaceDefStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 3, Column: 1}
		id := &ast.InterfaceDef{Pos: pos, Name: "Drawable"}
		var s ast.Statement = id
		assert.Equal(t, pos, s.Position())
	})
}

func TestEnumDefStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 4, Column: 1}
		ed := &ast.EnumDef{Pos: pos, Name: "Color"}
		var s ast.Statement = ed
		assert.Equal(t, pos, s.Position())
	})
}

func TestRelationshipStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 5, Column: 1}
		r := &ast.Relationship{Pos: pos, Left: "A", Right: "B"}
		var s ast.Statement = r
		assert.Equal(t, pos, s.Position())
	})
}

func TestPackageStatement(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsStatement", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 6, Column: 1}
		p := &ast.Package{Pos: pos, Name: "com.example"}
		var s ast.Statement = p
		assert.Equal(t, pos, s.Position())
	})
}

func TestFieldMember(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsMember", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 7, Column: 3}
		f := &ast.Field{Pos: pos, Name: "name", Type: "String"}
		var m ast.Member = f
		assert.Equal(t, pos, m.Position())
	})
}

func TestMethodMember(t *testing.T) {
	t.Parallel()
	t.Run("ImplementsMember", func(t *testing.T) {
		t.Parallel()
		pos := lexer.Pos{Line: 8, Column: 3}
		m := &ast.Method{Pos: pos, Name: "speak"}
		var member ast.Member = m
		assert.Equal(t, pos, member.Position())
	})
}

func TestRelationshipTypeConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.RelationshipType(0), ast.RelAssociation)
		assert.Equal(t, ast.RelationshipType(1), ast.RelDependency)
		assert.Equal(t, ast.RelationshipType(2), ast.RelInheritance)
		assert.Equal(t, ast.RelationshipType(3), ast.RelRealization)
		assert.Equal(t, ast.RelationshipType(4), ast.RelComposition)
		assert.Equal(t, ast.RelationshipType(5), ast.RelAggregation)
	})
}

func TestArrowDirectionConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.ArrowDirection(0), ast.ArrowLeft)
		assert.Equal(t, ast.ArrowDirection(1), ast.ArrowRight)
		assert.Equal(t, ast.ArrowDirection(2), ast.ArrowBoth)
		assert.Equal(t, ast.ArrowDirection(3), ast.ArrowNone)
	})
}

func TestVisibilityConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.Visibility(0), ast.VisibilityNone)
		assert.Equal(t, ast.Visibility(1), ast.VisibilityPublic)
		assert.Equal(t, ast.Visibility(2), ast.VisibilityPrivate)
		assert.Equal(t, ast.Visibility(3), ast.VisibilityProtected)
		assert.Equal(t, ast.Visibility(4), ast.VisibilityPackage)
	})
}

func TestModifierConstants(t *testing.T) {
	t.Parallel()
	t.Run("Values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ast.Modifier(0), ast.ModifierNone)
		assert.Equal(t, ast.Modifier(1), ast.ModifierStatic)
		assert.Equal(t, ast.Modifier(2), ast.ModifierField)
		assert.Equal(t, ast.Modifier(3), ast.ModifierMethod)
	})
}
