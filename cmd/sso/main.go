package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/TemaStatham/sso/internal/app"
	"github.com/TemaStatham/sso/internal/config"
	authservice "github.com/TemaStatham/sso/internal/services/auth"
	authstorage "github.com/TemaStatham/sso/internal/storage/sqlite"
)

// go run ./cmd/sso/main.go --config="./config/config.yaml"

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.Any("cfg", cfg))
	storage, err := authstorage.New(cfg.StoragePaths)
	if err != nil {
		panic(err)
	}
	service := authservice.New(log, storage,storage,storage,cfg.TokenTTL)
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePaths, cfg.TokenTTL, service)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("get signal stop", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
