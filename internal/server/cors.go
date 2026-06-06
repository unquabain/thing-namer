package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ReferrerCORSMiddleware sets Access-Control-Allow-Origin based on the
// scheme+host of the Referer header. Preflight OPTIONS requests are answered
// with 204; non-OPTIONS requests always run the wrapped handler even when
// the Referer is missing or malformed (X-Error is set in those cases).
func ReferrerCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		referrer := r.Referer()
		if referrer == "" {
			w.Header().Add("X-Error", "missing referrer")
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(referrer, "http://") && !strings.HasPrefix(referrer, "https://") {
			w.Header().Add("X-Error", fmt.Sprintf("invalid referrer: %s", referrer))
			next.ServeHTTP(w, r)
			return
		}
		rurl, err := url.Parse(referrer)
		if err != nil {
			w.Header().Add("X-Error", fmt.Sprintf("invalid referrer: %s", err.Error()))
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", fmt.Sprintf("%s://%s", rurl.Scheme, rurl.Host))
		next.ServeHTTP(w, r)
	})
}
