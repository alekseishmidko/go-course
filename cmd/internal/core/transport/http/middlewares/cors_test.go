package core_http_middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	for _, tc := range []struct {
		name, method, origin string
		allowed              bool
	}{
		{"allowed origin", "GET", "http://localhost:5050", true},
		{"disallowed origin", "GET", "http://localhost:3000", false},
		{"no origin", "GET", "", false},
		{"allowed preflight", "OPTIONS", "http://localhost:5050", true},
		{"disallowed preflight", "OPTIONS", "https://example.com", false},
		{"ordinary options", "OPTIONS", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			h := CORS()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusAccepted)
			}))
			r := httptest.NewRequest(tc.method, "/api/v1/tasks", nil)
			r.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			wantHeaders := map[string]string{
				"Access-Control-Allow-Origin":  tc.origin,
				"Access-Control-Allow-Methods": "GET, POST, PATCH, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			}
			for key, want := range wantHeaders {
				if !tc.allowed {
					want = ""
				}
				if got := w.Header().Get(key); got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			if tc.method == http.MethodOptions {
				if called || w.Code != http.StatusOK || w.Body.Len() != 0 {
					t.Fatalf("unexpected OPTIONS response: called=%v status=%d body=%q", called, w.Code, w.Body.String())
				}
			} else if !called || w.Code != http.StatusAccepted {
				t.Fatalf("handler response lost: called=%v status=%d", called, w.Code)
			}
		})
	}
}
