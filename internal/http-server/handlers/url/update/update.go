package update

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	urlservice "url-shortener/internal/service/url"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	NewURL string `json:"url"`
	Alias  string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	CountUpdated int64 `json:"countUpdated"`
}

type UpdateURL interface {
	UpdateURL(alias, newURL string) (int64, error)
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

		countUpdated, err := updateURL.UpdateURL(req.Alias, req.NewURL)
		if errors.Is(err, urlservice.ErrInvalidURL) {
			render.JSON(w, r, resp.Error("invalid url"))
			log.Info("invalid url", slog.String("url", req.NewURL))
			return
		}
		if errors.Is(err, urlservice.ErrURLNotFound) {
			render.JSON(w, r, resp.Error("url not found"))
			log.Info("url not found", slog.String("alias", req.Alias))
			return
		}
		if err != nil {
			log.Error("failed to update url", slog.String("alias", req.Alias), sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated", "alias", req.Alias, "countUpdated", countUpdated)

		responseOK(w, r, countUpdated)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, countUpdated int64) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		CountUpdated: countUpdated,
	})
}
