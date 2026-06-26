package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	core_postgres_pool "github.com/IvanJSBog/goland-todo-app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/IvanJSBog/goland-todo-app/internal/features/users/service"
	users_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("Failed to initialize logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	// initializing features
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	err = httpServer.Start(ctx)
	if err != nil {
		logger.Error("Http server start error", zap.Error(err))
	}
}
