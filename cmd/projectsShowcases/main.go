package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"projectsShowcase/internal/config"
	"projectsShowcase/internal/http-server/handlers/application/getAll"
	"projectsShowcase/internal/http-server/handlers/application/getApproved"
	"projectsShowcase/internal/http-server/handlers/application/remove"
	"projectsShowcase/internal/http-server/handlers/application/save"
	"projectsShowcase/internal/http-server/handlers/application/updateStatus"
	"projectsShowcase/internal/http-server/middleware/logger"
	"projectsShowcase/internal/lib/logger/sl"
	"projectsShowcase/internal/storage/sqlite"
	"syscall"
	"time"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.New(log))

	router.Post("/applications", save.New(log, storage))
	router.Get("/applications/approved", getApproved.New(log, storage))

	router.Route("/admin", func(r chi.Router) {
		r.Use(middleware.BasicAuth("projects-showcase", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Get("/applications", getAll.New(log, storage))
		r.Patch("/applications/{id}", updateStatus.New(log, storage))
		r.Delete("/applications/{id}", remove.New(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
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
