package main

import (
	"github.com/DaniilKalts/url-shortener/internal/storage/sqlite"
	"log/slog"
	"os"

	"github.com/DaniilKalts/url-shortener/internal/config"
	mySlog "github.com/DaniilKalts/url-shortener/lib/logger/slog"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("Debug logger", slog.Any("env", cfg))

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		logger.Error("Error opening storage", mySlog.Err)
		os.Exit(1)
	}

	_ = storage

	// TO-DO: init router (chi)

	// TO-DO: run server (http)
}

type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "dev"
	EnvProd  Environment = "prod"
)

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "local":
		opts := slog.HandlerOptions{Level: slog.LevelDebug}
		logger = slog.New(slog.NewTextHandler(os.Stdout, &opts))
	case "dev":
		opts := slog.HandlerOptions{Level: slog.LevelDebug}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	case "prod":
		opts := slog.HandlerOptions{Level: slog.LevelInfo}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	}

	return logger
}
