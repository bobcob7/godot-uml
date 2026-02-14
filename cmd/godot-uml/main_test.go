package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "godot-uml")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(".", "")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))
	return bin
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "input.puml")
	require.NoError(t, os.WriteFile(f, []byte(content), 0o644))
	return f
}

const validClass = "@startuml\nclass Foo {\n+name : String\n}\n@enduml"

func TestParseRenderArgs(t *testing.T) {
	t.Parallel()
	t.Run("FileOnly", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"input.puml"})
		assert.Equal(t, "", out)
		assert.Equal(t, "input.puml", in)
	})
	t.Run("FileWithOutputAfter", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"input.puml", "-o", "out.svg"})
		assert.Equal(t, "out.svg", out)
		assert.Equal(t, "input.puml", in)
	})
	t.Run("OutputBeforeFile", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"-o", "out.svg", "input.puml"})
		assert.Equal(t, "out.svg", out)
		assert.Equal(t, "input.puml", in)
	})
	t.Run("Stdin", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"-"})
		assert.Equal(t, "", out)
		assert.Equal(t, "-", in)
	})
	t.Run("StdinWithOutput", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"-", "-o", "out.svg"})
		assert.Equal(t, "out.svg", out)
		assert.Equal(t, "-", in)
	})
	t.Run("Help", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{"--help"})
		assert.Equal(t, "", out)
		assert.Equal(t, "", in)
	})
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		out, in := parseRenderArgs([]string{})
		assert.Equal(t, "", out)
		assert.Equal(t, "", in)
	})
}

func TestCmdRender(t *testing.T) {
	t.Parallel()
	t.Run("FileToFile", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, validClass)
		output := filepath.Join(t.TempDir(), "out.svg")
		code := cmdRender([]string{input, "-o", output})
		assert.Equal(t, exitSuccess, code)
		data, err := os.ReadFile(output)
		require.NoError(t, err)
		assert.Contains(t, string(data), "<svg")
		assert.Contains(t, string(data), "Foo")
	})
	t.Run("OutputBeforeFile", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, validClass)
		output := filepath.Join(t.TempDir(), "out.svg")
		code := cmdRender([]string{"-o", output, input})
		assert.Equal(t, exitSuccess, code)
		data, err := os.ReadFile(output)
		require.NoError(t, err)
		assert.Contains(t, string(data), "<svg")
	})
	t.Run("MissingFile", func(t *testing.T) {
		t.Parallel()
		code := cmdRender([]string{"/nonexistent/file.puml"})
		assert.Equal(t, exitSystem, code)
	})
	t.Run("NoArgs", func(t *testing.T) {
		t.Parallel()
		code := cmdRender([]string{})
		assert.Equal(t, exitSystem, code)
	})
	t.Run("InvalidDiagram", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, "not a diagram")
		output := filepath.Join(t.TempDir(), "out.svg")
		code := cmdRender([]string{input, "-o", output})
		assert.Equal(t, exitValidation, code)
	})
}

func TestCmdValidate(t *testing.T) {
	t.Parallel()
	t.Run("ValidFile", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, validClass)
		code := cmdValidate([]string{input})
		assert.Equal(t, exitSuccess, code)
	})
	t.Run("InvalidFile", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, "not a diagram")
		code := cmdValidate([]string{input})
		assert.Equal(t, exitValidation, code)
	})
	t.Run("MissingFile", func(t *testing.T) {
		t.Parallel()
		code := cmdValidate([]string{"/nonexistent/file.puml"})
		assert.Equal(t, exitSystem, code)
	})
	t.Run("NoArgs", func(t *testing.T) {
		t.Parallel()
		code := cmdValidate([]string{})
		assert.Equal(t, exitSystem, code)
	})
}

func TestBinary(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	t.Run("Version", func(t *testing.T) {
		t.Parallel()
		out, err := exec.Command(bin, "version").CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(out), "godot-uml")
	})
	t.Run("Help", func(t *testing.T) {
		t.Parallel()
		cmd := exec.Command(bin, "help")
		out, _ := cmd.CombinedOutput()
		assert.Contains(t, string(out), "Usage:")
		assert.Contains(t, string(out), "render")
		assert.Contains(t, string(out), "validate")
		assert.Contains(t, string(out), "serve")
	})
	t.Run("NoArgs", func(t *testing.T) {
		t.Parallel()
		cmd := exec.Command(bin)
		out, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.Contains(t, string(out), "Usage:")
	})
	t.Run("UnknownCommand", func(t *testing.T) {
		t.Parallel()
		cmd := exec.Command(bin, "bogus")
		out, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.Contains(t, string(out), "unknown command")
	})
	t.Run("RenderFile", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, validClass)
		output := filepath.Join(t.TempDir(), "out.svg")
		out, err := exec.Command(bin, "render", input, "-o", output).CombinedOutput()
		require.NoError(t, err, "render failed: %s", string(out))
		data, err := os.ReadFile(output)
		require.NoError(t, err)
		assert.Contains(t, string(data), "<svg")
	})
	t.Run("RenderStdin", func(t *testing.T) {
		t.Parallel()
		cmd := exec.Command(bin, "render", "-")
		cmd.Stdin = os.Stdin
		f := writeTempFile(t, validClass)
		stdin, err := os.Open(f)
		require.NoError(t, err)
		defer func() { _ = stdin.Close() }()
		cmd.Stdin = stdin
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "render stdin failed: %s", string(out))
		assert.Contains(t, string(out), "<svg")
	})
	t.Run("ValidateValid", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, validClass)
		out, err := exec.Command(bin, "validate", input).CombinedOutput()
		require.NoError(t, err)
		assert.Contains(t, string(out), "OK")
	})
	t.Run("ValidateInvalid", func(t *testing.T) {
		t.Parallel()
		input := writeTempFile(t, "not a diagram")
		cmd := exec.Command(bin, "validate", input)
		out, err := cmd.CombinedOutput()
		assert.Error(t, err)
		assert.Contains(t, string(out), "expected @startuml")
	})
}
