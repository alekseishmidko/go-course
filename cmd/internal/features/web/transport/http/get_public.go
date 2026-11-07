package web_transport_http

import (
	"bytes"
	"errors"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
	"net/http"
)

func (h *WebHTTPHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, err := h.webService.GetFile(r.Context(), r.URL.Path)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to get public file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, file.Name, file.ModTime, bytes.NewReader(file.Content))
}
