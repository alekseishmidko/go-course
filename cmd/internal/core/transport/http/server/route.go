package core_http_server

import (
	"net/http"

	core_http_middlewares "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/middlewares"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middlewares.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middlewares.ChainMiddleware(r.Handler, r.Middleware...)
}
