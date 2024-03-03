package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/TemaStatham/sso/internal/app/grpc"
	authgrpc "github.com/TemaStatham/sso/internal/grpc/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration, authService authgrpc.Auth) *App {
	// TODO: инициализировать хранилище

	// TODO: инициализировать auth service

	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
