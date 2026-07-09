package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/IvanJSBog/goland-todo-app/docs"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_http_middleware "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(config Config, log *core_logger.Logger, middleware ...core_http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		for _, route := range router.routes {
			pattern := fmt.Sprintf(
				"%s /api/%s%s",
				route.Method,
				router.apiVersion,
				route.Path,
			)

			s.mux.Handle(pattern, route.Handler)
		}
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)
	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Start(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("starting http server on ", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("http server start: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutting down http server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("http server shutdown: %w", err)
		}
	}

	return nil
}
