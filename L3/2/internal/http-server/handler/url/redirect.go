package url

import (
	"errors"
	"log/slog"
	"net/http"

	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type URLRedirector interface {
	GetURL(alias string) (string, error)
	UpdateURLStats(alias string, userAgent string) error
}

func Redirect(logger *slog.Logger, urlRedirector URLRedirector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.Redirect"
		logger = logger.With(slog.String("op", op))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "invalid request",
			})
			return
		}

		resURL, err := urlRedirector.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", "alias", alias)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "not found",
			})
			return
		}
		if err != nil {
			logger.Error("failed to get url", slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "internal error",
			})
			return
		}
		logger.Info("got url", slog.String("url", resURL))

		if err := urlRedirector.UpdateURLStats(alias, r.UserAgent()); err != nil {
			logger.Error("failed to update stats", slog.Any("err", err))
		}

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
