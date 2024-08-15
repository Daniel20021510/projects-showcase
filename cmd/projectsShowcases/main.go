package main

import (
	"log/slog"
	"os"
	"projectsShowcase/internal/config"
	"projectsShowcase/internal/lib/logger/sl"
	"projectsShowcase/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	_ = storage
}

// setupLogger returns a logger based on the environment.
//
// The function takes an environment string as input and returns a pointer to a slog.Logger.
// The logger is configured based on the environment:
//
// - If the environment is envLocal, a text handler with debug level is used.
//
// - If the environment is envDev, a JSON handler with debug level is used.
//
// - If the environment is envProd, a JSON handler with info level is used.
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
