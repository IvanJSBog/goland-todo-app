package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}
