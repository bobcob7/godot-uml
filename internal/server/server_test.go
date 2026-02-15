package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bobcob7/go-uml/internal/encoding"
	"github.com/bobcob7/go-uml/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer() http.Handler {
	return server.New(server.DefaultConfig()).Handler()
}

func TestServer(t *testing.T) {
	t.Parallel()
	t.Run("PostRenderClassDiagram", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		body := "@startuml\nclass Foo {\n+name : String\n}\n@enduml"
		req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "image/svg+xml", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "<svg")
		assert.Contains(t, rec.Body.String(), "Foo")
	})
	t.Run("PostRenderSequenceDiagram", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		body := "@startuml\nparticipant Alice\nparticipant Bob\nAlice -> Bob : hello\n@enduml"
		req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<svg")
		assert.Contains(t, rec.Body.String(), "Alice")
	})
	t.Run("PostRenderInvalidDiagram", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		body := "not a diagram"
		req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), `"errors"`)
		assert.Contains(t, rec.Body.String(), `"line"`)
		assert.Contains(t, rec.Body.String(), `"message"`)
	})
	t.Run("PostRenderEmptyBody", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(""))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		// Empty input produces a valid empty SVG (parser handles it gracefully)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<svg")
	})
	t.Run("GetSVGEncoded", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		text := "@startuml\nclass Foo\n@enduml"
		encoded, err := encoding.Encode(text)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/svg/"+encoded, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "image/svg+xml", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "<svg")
		assert.Contains(t, rec.Body.String(), "Foo")
	})
	t.Run("GetSVGInvalidEncoding", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		req := httptest.NewRequest(http.MethodGet, "/svg/!!!!invalid", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("GetEditor", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
		body, _ := io.ReadAll(rec.Body)
		assert.Contains(t, string(body), "go-uml")
		assert.Contains(t, string(body), "<textarea")
		assert.Contains(t, string(body), "/render")
	})
	t.Run("GetEditorNotFoundPath", func(t *testing.T) {
		t.Parallel()
		handler := newTestServer()
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := server.DefaultConfig()
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Greater(t, cfg.ReadTimeout.Seconds(), 0.0)
	assert.Greater(t, cfg.WriteTimeout.Seconds(), 0.0)
}
