package save

import (
	"errors"
	"github.com/DaniilKalts/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/DaniilKalts/url-shortener/lib/api/response"
	mySlog "github.com/DaniilKalts/url-shortener/lib/logger/slog"
	"github.com/DaniilKalts/url-shortener/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

// To-Do: Move to config
const aliasLength = 4

type URLSaver interface {
	SaveURL(alias string, url string) (int, error)
}

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operation := "http-server.handlers.save.New"

		logger = logger.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request", mySlog.Err(err))

			render.JSON(w, r, response.Error("Failed to decode request"))

			return
		}

		logger.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			logger.Error("Failed to validate request", mySlog.Err(err))

			render.JSON(w, r, response.ValidationError(validationErrs))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(alias, req.URL)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Error("URL already exists", mySlog.Err(err))

			render.JSON(w, r, response.Error("URL already exists"))

			return
		}
		if err != nil {
			logger.Error("Failed to save url", mySlog.Err(err))

			render.JSON(w, r, response.Error("Failed to save url"))

			return
		}

		logger.Info("Url added", slog.Int("id", id))

		render.JSON(
			w, r, Response{
				Response: response.OK(),
				Alias:    alias,
			},
		)
	}
}
