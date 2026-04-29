package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_middlewares "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/middlewares"
	core_http_server "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/server"
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
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(nil)

	userRoutes := usersTransportHttp.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger, core_http_middlewares.RequestID(),
		core_http_middlewares.Logger(logger),
		core_http_middlewares.Panic(),
		core_http_middlewares.Trace())

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run Error", zap.Error(err))
	}
}
