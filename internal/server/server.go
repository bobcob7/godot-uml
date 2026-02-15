// Package server provides the HTTP server and live editor for go-uml.
package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bobcob7/go-uml/internal/encoding"
	"github.com/bobcob7/go-uml/pkg/gouml"
)

//go:embed static/*
var staticFS embed.FS

// Config holds server configuration.
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Server is the HTTP server for go-uml.
type Server struct {
	config Config
	mux    *http.ServeMux
}

// New creates a new Server with the given config.
func New(cfg Config) *Server {
	s := &Server{config: cfg, mux: http.NewServeMux()}
	s.mux.HandleFunc("POST /render", s.handleRender)
	s.mux.HandleFunc("GET /svg/{encoded...}", s.handleSVG)
	s.mux.HandleFunc("GET /", s.handleEditor)
	return s
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}
	return srv.ListenAndServe()
}

// Handler returns the HTTP handler for testing.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	errs := gouml.Validate(strings.NewReader(string(body)))
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := errorResponse{Errors: make([]errorDetail, len(errs))}
		for i, e := range errs {
			resp.Errors[i] = errorDetail{Line: e.Line, Column: e.Column, Message: e.Message}
		}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	if err := gouml.Render(strings.NewReader(string(body)), w); err != nil {
		http.Error(w, fmt.Sprintf("render error: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSVG(w http.ResponseWriter, r *http.Request) {
	encoded := r.PathValue("encoded")
	if encoded == "" {
		http.Error(w, "missing encoded diagram", http.StatusBadRequest)
		return
	}
	text, err := encoding.Decode(encoded)
	if err != nil {
		http.Error(w, fmt.Sprintf("decode error: %s", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	if err := gouml.Render(strings.NewReader(text), w); err != nil {
		http.Error(w, fmt.Sprintf("render error: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleEditor(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "editor not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

type errorResponse struct {
	Errors []errorDetail `json:"errors"`
}

type errorDetail struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}
