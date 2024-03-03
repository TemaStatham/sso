package grpc

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/TemaStatham/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

// App это grpc server app
type App struct {
	log     *slog.Logger
	gRPCSrv *grpc.Server
	port    int
}

// New создает новый grpc server app
func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:     log,
		gRPCSrv: gRPCServer,
		port:    port,
	}
}

// MustRun обертка для метода Run
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run стартует сервер
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	log.Info("starting grpc server")

	if err := a.gRPCSrv.Serve(l); err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}

// Stop останавливает приложение
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op), slog.Any("stopping grpc server", slog.Int("port", a.port)))

	a.gRPCSrv.GracefulStop()
}
