package web_transport_http

import (
	"context"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"

	core_http_server "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetFile(ctx context.Context, urlPath string) (domain.WebFile, error)
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{webService: webService}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	// An unrestricted pattern avoids conflicts with the outer /api/v1/ router.
	// Supported methods are checked by GetPublic instead.
	return []core_http_server.Route{{Path: "/", Handler: h.GetPublic}}
}
