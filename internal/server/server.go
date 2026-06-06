package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig"

	"github.com/Unquabain/thing-namer/internal/namer"
)

// Server wires a Namer to HTTP handlers. The zero value is not usable;
// construct with New.
type Server struct {
	namer   *namer.Namer
	index   *template.Template
	goTmplt *template.Template
}

// New returns a Server using the provided Namer.
func New(n *namer.Namer) *Server {
	return &Server{
		namer:   n,
		index:   template.Must(template.New("index").Funcs(sprig.TxtFuncMap()).Parse(indexRaw)),
		goTmplt: template.Must(template.New("go").Funcs(sprig.TxtFuncMap()).Parse(goRaw)),
	}
}

// Handler returns the fully-wrapped http.Handler (CORS middleware + routes).
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serve)
	return ReferrerCORSMiddleware(mux)
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	switch {
	case requestIsJSON(r):
		s.renderJSON(w)
	case requestIsGo(r):
		s.renderGo(w, r)
	default:
		s.renderHTML(w)
	}
}

func requestIsJSON(r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, ".json") {
		return true
	}
	return r.Header.Get("Accept") == "application/json"
}

func requestIsGo(r *http.Request) bool {
	return strings.HasSuffix(r.URL.Path, ".go")
}

func (s *Server) renderJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	ctx := s.namer.NewContext()
	if err := json.NewEncoder(w).Encode(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	ctx := s.namer.NewThemedContext()
	if err := s.index.Execute(w, ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderGo(w http.ResponseWriter, r *http.Request) {
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	if xfp := r.Header.Get("X-Forwarded-Proto"); xfp != "" {
		proto = xfp
	}
	ctx := struct{ Server string }{fmt.Sprintf("%s://%s", proto, r.Host)}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", `attachment; filename="wizardbacon.go"`)
	if err := s.goTmplt.Execute(w, ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
