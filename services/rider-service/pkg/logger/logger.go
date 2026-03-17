package logger

import (
	"log/slog"
	"os"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
)

func New(cfg config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.App.Environment != "production" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler).With("service", cfg.App.Name, "env", cfg.App.Environment)
}
