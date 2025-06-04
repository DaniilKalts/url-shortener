package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/DaniilKalts/url-shortener/internal/storage"
	"github.com/DaniilKalts/url-shortener/lib/api/response"
	mySlog "github.com/DaniilKalts/url-shortener/lib/logger/slog"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operation := "http-server.handlers.delete.New"
		logger = logger.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("Alias is empty")
			render.JSON(w, r, response.Error("Alias is empty"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("URL not found", slog.String("alias", alias))
			render.JSON(w, r, response.Error("URL not found"))
			return
		}
		if err != nil {
			logger.Error("Failed to delete URL", mySlog.Err(err))
			render.JSON(w, r, response.Error("Failed to delete URL"))
			return
		}

		logger.Info("URL deleted", slog.String("alias", alias))
		render.JSON(w, r, response.OK())
	}
}
