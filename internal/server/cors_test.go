package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCORSOptionsWithOriginSetsHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://foo.example")
	rr := httptest.NewRecorder()
	ReferrerCORSMiddleware(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rr.Code)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "https://foo.example" {
		t.Fatalf("ACAO = %q", got)
	}
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("ACAM not set")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatal("ACAH not set")
	}
}

func TestCORSOptionsWithoutOriginReturns204(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rr := httptest.NewRecorder()
	ReferrerCORSMiddleware(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rr.Code)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("ACAO unexpectedly set: %q", got)
	}
}

func TestCORSGetMissingReferer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	ReferrerCORSMiddleware(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (handler ran)", rr.Code)
	}
	if got := rr.Header().Get("X-Error"); got != "missing referrer" {
		t.Fatalf("X-Error = %q", got)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("ACAO unexpectedly set: %q", got)
	}
}

func TestCORSGetInvalidScheme(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Referer", "ftp://bad.example/")
	rr := httptest.NewRecorder()
	ReferrerCORSMiddleware(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if got := rr.Header().Get("X-Error"); got == "" {
		t.Fatal("X-Error not set for invalid referrer")
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("ACAO unexpectedly set: %q", got)
	}
}

func TestCORSGetValidReferer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Referer", "https://example.com/some/path")
	rr := httptest.NewRecorder()
	ReferrerCORSMiddleware(okHandler()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("ACAO = %q, want https://example.com", got)
	}
}
