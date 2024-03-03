package main

import (
	"log/slog"
	"os"

	"github.com/TemaStatham/sso/internal/config"
)

// go run ./cmd/sso/main.go --config="./config/config.yaml"

const (
	envLocal="local"
	envDev="dev"
	envProd="prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.Any("cfg", cfg))
	log.Debug("debug message")
	log.Error("error mesage")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch (env) {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}