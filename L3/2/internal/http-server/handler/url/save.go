package url

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"url-shortener/internal/storage"
	"url-shortener/lib/random"

	"github.com/go-chi/render"
)

const aliasLength = 5

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

type Request struct {
	URLToSave string `json:"url"`
	Alias     string `json:"alias"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Alias  string `json:"alias"`
}

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.save.New"
		logger = logger.With(slog.String("op", op))

		if r.Method != http.MethodPost {
			logger.Info("invalid method", slog.String("method", r.Method))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "only POST allowed",
			})
			return
		}

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to decode request body",
			})
			return
		}

		// Валидация: URL обязателен
		if req.URLToSave == "" {
			logger.Info("empty URL provided")
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "url is required",
			})
			return
		}

		logger.Info("request decoded successfully",
			slog.String("alias", req.Alias),
			// Не логируем полный URL! Безопасность!
			slog.String("url_length", strconv.Itoa(len(req.URLToSave))),
		)

		if req.Alias == "" {
			req.Alias = random.Alias(aliasLength)
		}

		_, err = urlSaver.SaveURL(req.URLToSave, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Info("url already exists", slog.String("alias", req.Alias), slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  err.Error(),
			})
			return
		}
		if err != nil {
			logger.Error("failed to save url", slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "internal server error", // лучше скрывать детали
			})
			return
		}

		logger.Info("url saved successfully", slog.String("alias", req.Alias))
		render.JSON(w, r, Response{
			Status: "success",
			Error:  "",
			Alias:  req.Alias,
		})
	}
}
