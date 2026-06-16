package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_postgres_pool "github.com/alekseishmidko/go-course/cmd/internal/core/repository/postgres/pool"
	core_http_middlewares "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/middlewares"
	core_http_server "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/server"
	users_postgres_repository "github.com/alekseishmidko/go-course/cmd/internal/features/users/repository/postgres"
	users_service "github.com/alekseishmidko/go-course/cmd/internal/features/users/service"
	users_transport_http "github.com/alekseishmidko/go-course/cmd/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Println("Failed to init app logger", err)
		os.Exit(1)
	}
	defer logger.Close()
	logger.Debug("Starting application")

	logger.Debug("Initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("Failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()
	// USERS
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing https server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger, core_http_middlewares.RequestID(),
		core_http_middlewares.Logger(logger),
		core_http_middlewares.Panic(),
		core_http_middlewares.Trace())

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run Error", zap.Error(err))
	}
}
