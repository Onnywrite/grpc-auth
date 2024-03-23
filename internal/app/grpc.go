package app

import (
	"fmt"
	"log/slog"
	"net"

	grpcauth "github.com/Onnywrite/grpc-auth/internal/grpc/auth"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log    *slog.Logger
	server *grpc.Server
	port   string
}

func NewGRPC(logger *slog.Logger, service grpcauth.AuthService, port int) *GRPCApp {
	s := grpc.NewServer()

	// add middlewares if possible

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
			slog.String("err", err.Error()),
		)
		return err
	}

	a.log.Info("gRPC started",
		slog.String("port", a.port),
	)

	if err := a.server.Serve(lis); err != nil {
		a.log.Error("error while starting gRPC",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		return err
	}

	return nil
}

func (a *GRPCApp) Stop() {
	const op = "GRPCApp.Stop"

	a.log.Info("stopping gRPC",
		slog.String("port", a.port),
	)

	a.server.GracefulStop()
}
