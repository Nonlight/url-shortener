package update

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	NewURL string `json:"url" validate:"required,url"`
	Alias  string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	CountUpdated int64 `json:"counteUpdate"`
}

type UpdateURL interface {
	UpdateURL(urlToUpdate string, newURL string) (int64, error)
}

func New(log *slog.Logger, updateURL UpdateURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.update.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body", slog.Any("url", req))

		if err := validator.New().Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(validationErrors))
			render.JSON(w, r, resp.ValidationError(validationErrors))
			return
		}

		alias := req.Alias
		newURL := req.NewURL

		countUpdated, err := updateURL.UpdateURL(alias, newURL)
		if errors.Is(err, storage.ErrURLNotFound) {
			render.JSON(w, r, resp.Error("url not found"))
			log.Info("url not found", "alias", alias)
			return
		}
		if err != nil {
			log.Error("failed to update url", "alias", alias, sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated", "alias", alias, "countUpdated", countUpdated)

		responseOK(w, r, countUpdated)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, countUpdated int64) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		CountUpdated: countUpdated,
	})
}
