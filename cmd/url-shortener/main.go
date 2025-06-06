package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DaniilKalts/url-shortener/internal/config"
	"github.com/DaniilKalts/url-shortener/internal/http-server/handlers/redirect"
	"github.com/DaniilKalts/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/DaniilKalts/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/DaniilKalts/url-shortener/internal/http-server/middlewares/logger"
	"github.com/DaniilKalts/url-shortener/internal/storage/sqlite"
	"github.com/DaniilKalts/url-shortener/lib/logger/handlers/slogpretty"
	mySlog "github.com/DaniilKalts/url-shortener/lib/logger/slog"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("Debug logger", slog.Any("env", cfg))

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		logger.Error("Error opening storage", mySlog.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route(
		"/url", func(r chi.Router) {
			r.Use(
				middleware.BasicAuth(
					"url-shortener", map[string]string{
						cfg.HTTPServer.User: cfg.HTTPServer.Password,
					},
				),
			)

			r.Post("/", save.New(logger, storage))
			r.Delete("/{alias}", delete.New(logger, storage))
		},
	)
	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("Starting server...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Error starting server", mySlog.Err(err))
	}
}

func setupLogger(env config.Environment) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = setupPrettySlog()
	case config.EnvDev:
		opts := slog.HandlerOptions{Level: slog.LevelDebug}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	case config.EnvProd:
		opts := slog.HandlerOptions{Level: slog.LevelInfo}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug}}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
