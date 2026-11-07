package core_http_server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootRoutesDoNotConflictWithAPI(t *testing.T) {
	s := &HTTPServer{mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	s.RegisterSwagger()
	s.RegisterRoutes(Route{Path: "/", Handler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}})
	for path, status := range map[string]int{"/api/v1/tasks": http.StatusNoContent, "/swagger/doc.json": http.StatusOK} {
		w := httptest.NewRecorder()
		s.mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != status {
			t.Errorf("%s: status = %d, want %d", path, w.Code, status)
		}
	}
}
