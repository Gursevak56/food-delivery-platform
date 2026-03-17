package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/app"
	bootstrap "github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/bootstrap"
	pkglogger "github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	logger := pkglogger.New(cfg)

	container, err := app.NewContainer(cfg, logger)
	if err != nil {
		logger.Error("failed to bootstrap rider service", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	server := bootstrap.NewHTTPServer(cfg, logger, bootstrap.NewRouter(container))

	go func() {
		logger.Info("starting rider service", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down rider service")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	_ = server.Shutdown(ctx)
}
