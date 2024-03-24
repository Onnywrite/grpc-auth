package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/storage"
)

type App struct {
	log  *slog.Logger
	grpc *GRPCApp
	db   *storage.Pg
}

func New(logger *slog.Logger, conn string, tokenTTL time.Duration, grpcPort int, grpcTimeout time.Duration) *App {
	const op = "app.New"

	db, err := storage.NewPg(conn)
	if err != nil {
		logger.Error("could not connect to pg database",
			slog.String("op", op),
			slog.String("err", err.Error()))
		os.Exit(1)
	}

	//authService := auth.New(...)

	app := NewGRPC(logger, nil /*authService*/, grpcPort)

	return &App{
		grpc: app,
		log:  logger,
		db:   db,
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

	if err := a.db.Disconnect(); err != nil {
		a.log.Error("could not disconnect from pg database",
			slog.String("op", op),
			slog.String("err", err.Error()))
		return
	}
	a.log.Info("successfully disconnected from pg database", slog.String("op", op))
}
