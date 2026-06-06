package server

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Unquabain/thing-namer/internal/namer"
)

func loadTestNamer(t *testing.T) *namer.Namer {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// internal/server -> repo root
	root := filepath.Join(wd, "..", "..")
	data, err := os.ReadFile(filepath.Join(root, "data", "words.yaml"))
	if err != nil {
		t.Fatalf("read words.yaml: %v", err)
	}
	n, err := namer.New(data, namer.WithRand(rand.New(rand.NewSource(1))))
	if err != nil {
		t.Fatalf("namer.New: %v", err)
	}
	return n
}

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	return New(loadTestNamer(t)).Handler()
}

func TestRouteJSON(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api.json", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q", ct)
	}
	var body struct {
		ProjectName string `json:"projectName"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.ProjectName == "" {
		t.Fatal("ProjectName empty")
	}
}

func TestRouteHTML(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api.html", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/html" {
		t.Fatalf("Content-Type = %q", ct)
	}
	body, _ := io.ReadAll(rr.Body)
	if len(body) == 0 {
		t.Fatal("empty body")
	}
}

func TestRouteGo(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api.go", nil)
	req.Host = "example.test"
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("Content-Type = %q", ct)
	}
	body, _ := io.ReadAll(rr.Body)
	if _, err := parser.ParseFile(token.NewFileSet(), "client.go", body, parser.AllErrors); err != nil {
		t.Fatalf("rendered Go did not parse: %v\n---\n%s", err, body)
	}
}

func TestRouteRootRendersHTML(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/html" {
		t.Fatalf("Content-Type = %q", ct)
	}
}
