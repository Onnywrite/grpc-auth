package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Onnywrite/grpc-auth/internal/app"
	"github.com/Onnywrite/grpc-auth/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Environment)

	application := app.New(logger, cfg.Conn, cfg.TokenTTL, cfg.RefreshTokenTTL, cfg.GRPC.Port, cfg.GRPC.Timeout)
	go application.MustStart()

	shut := make(chan os.Signal, 1)
	signal.Notify(shut, syscall.SIGTERM, syscall.SIGINT)
	<-shut
	application.Stop()
	logger.Info("gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}
	return logger
}
