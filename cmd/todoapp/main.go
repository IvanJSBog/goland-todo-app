package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/IvanJSBog/goland-todo-app/docs"
	core_config "github.com/IvanJSBog/goland-todo-app/internal/core/config"
	core_logger "github.com/IvanJSBog/goland-todo-app/internal/core/logger"
	"github.com/IvanJSBog/goland-todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/IvanJSBog/goland-todo-app/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/statistics/repository/postgres"
	statistics_service "github.com/IvanJSBog/goland-todo-app/internal/features/statistics/service"
	statistics_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/IvanJSBog/goland-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/IvanJSBog/goland-todo-app/internal/features/users/service"
	users_transport_http "github.com/IvanJSBog/goland-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

// @title Golang Todo API
// @version 1.0
// @description Todo Application REST API scheme
// @host 127.0.0.1:5050
// @BasePath /api/v1
func main() {
	cnf := core_config.NewConfigMust()
	time.Local = cnf.TimeZone
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("Failed to initialize logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("time_zone", time.Local))

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

	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHttp := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHttp.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHttp.Routes()...)

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterSwagger()

	err = httpServer.Start(ctx)
	if err != nil {
		logger.Error("Http server start error", zap.Error(err))
	}
}
