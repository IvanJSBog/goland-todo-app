package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	"github.com/IvanJSBog/goland-todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/IvanJSBog/goland-todo-app/internal/features/users/service"
	users_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var timeZone = time.UTC

func main() {
	time.Local = timeZone
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("Failed to initialize logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("time_zone", timeZone))

	pool, err := core_pgx_pool.NewConnectionPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	// initializing features
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHttp.Routes()...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	err = httpServer.Start(ctx)
	if err != nil {
		logger.Error("Http server start error", zap.Error(err))
	}
}
