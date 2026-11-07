package web_transport_http

import (
	web_fs_repository "github.com/alekseishmidko/go-course/cmd/internal/features/web/repository/fs"
	web_service "github.com/alekseishmidko/go-course/cmd/internal/features/web/service"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetPublic(t *testing.T) {
	dir := t.TempDir()
	for name, content := range map[string]string{"index.html": "<!doctype html><h1>TodoApp</h1>", "app.js": "console.log('TodoApp');"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	h := NewWebHTTPHandler(web_service.NewWebService(web_fs_repository.NewWebRepository(dir)))
	mux := http.NewServeMux()
	for _, route := range h.Routes() {
		mux.Handle(route.Path, route.WithMiddleware())
	}
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/", 200}, {"GET", "/app.js", 200}, {"HEAD", "/", 200}, {"GET", "/missing", 404}, {"POST", "/", 405},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d", w.Code, tc.status)
			}
			if tc.method == "HEAD" && w.Body.Len() != 0 {
				t.Fatal("HEAD must have no body")
			}
			if tc.status == 405 && w.Header().Get("Allow") != "GET, HEAD" {
				t.Fatal("missing allowed methods")
			}
		})
	}
}
