package app

import (
	"log/slog"
	"time"
)

type App struct {
	log  *slog.Logger
	grpc *GRPCApp
}

func New(logger *slog.Logger, conn string, tokenTTL time.Duration, grpcPort int, grpcTimeout time.Duration) *App {
	// create new database for authService...

	//authService := auth.New(...)

	app := NewGRPC(logger, nil /*authService*/, grpcPort)

	return &App{
		grpc: app,
		log:  logger,
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	const op = "App.Start"

	a.log.Info("connecting to a database", slog.String("op", op))

	// ....

	a.log.Info("starting grpc", slog.String("op", op))
	return a.grpc.Start()
}

func (a *App) Stop() {
	const op = "App.Stop"

	a.log.Info("stopping application", slog.String("op", op))

	a.grpc.Stop()
}
