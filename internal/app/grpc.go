package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"
	"time"

	se "github.com/Onnywrite/grpc-auth/internal/lib/service-errors"
	"github.com/Onnywrite/grpc-auth/internal/transfer"
	grpcauth "github.com/Onnywrite/grpc-auth/internal/transfer/grpc/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCApp struct {
	log    *slog.Logger
	server *grpc.Server
	port   string
}

func NewGRPC(logger *slog.Logger, service transfer.AuthService, port int, timeout time.Duration) *GRPCApp {
	grpcLogger := logger.With(slog.String("op", "grpc"))

	s := grpc.NewServer(grpc.ConnectionTimeout(timeout), grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			grpcLogger.Error("recovered from panic", slog.Any("panic", p), slog.String("stack", string(debug.Stack())))

			return status.Errorf(codes.Internal, se.ErrPanicRecoveredGrpc.Error())
		})),
		logging.UnaryServerInterceptor(logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			grpcLogger.Log(ctx, slog.Level(lvl), msg, fields...)
		})),
	))

	grpcauth.Register(s, service)

	return &GRPCApp{
		log:    logger,
		server: s,
		port:   fmt.Sprintf(":%d", port),
	}
}

func (a *GRPCApp) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) Start() error {
	const op = "GRPCApp.Start"

	lis, err := net.Listen("tcp", a.port)
	if err != nil {
		a.log.Error("error while starting gRPC",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	a.log.Info("gRPC started",
		slog.String("port", a.port),
	)

	if err := a.server.Serve(lis); err != nil {
		a.log.Error("error while starting gRPC",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (a *GRPCApp) Stop() {
	const op = "GRPCApp.Stop"

	a.log.Info("stopping gRPC",
		slog.String("port", a.port),
		slog.String("op", op),
	)

	a.server.GracefulStop()
}
