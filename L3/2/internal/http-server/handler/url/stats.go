package url

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"

	"github.com/go-chi/render"
)

type Stat = storage.Stat

type StatsGetter interface {
	GetDailyStats() ([]storage.Stat, error)
	GetWeeklyStats() ([]storage.Stat, error)
	GetMonthlyStats() ([]storage.Stat, error)
}

func DailyStats(logger *slog.Logger, statsGetter StatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.stats.DailyStats"
		logger = logger.With(slog.String("op", op))

		stats, err := statsGetter.GetDailyStats()
		if err != nil {
			return
		}

		if stats == nil {
			stats = []storage.Stat{}
		}

		render.JSON(w, r, stats)
	}
}

func WeeklyStats(logger *slog.Logger, statsGetter StatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.stats.WeeklyStats"
		logger = logger.With(slog.String("op", op))

		stats, err := statsGetter.GetWeeklyStats()
		if err != nil {
			logger.Error("failed to get weekly stats", slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "internal error",
			})
			return
		}

		if stats == nil {
			stats = []storage.Stat{}
		}

		render.JSON(w, r, stats)
	}
}

func MonthlyStats(logger *slog.Logger, statsGetter StatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.stats.MonthlyStats"
		logger = logger.With(slog.String("op", op))

		stats, err := statsGetter.GetMonthlyStats()
		if err != nil {
			logger.Error("failed to get monthly stats", slog.Any("err", err))
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "internal error",
			})
			return
		}

		if stats == nil {
			stats = []storage.Stat{}
		}

		render.JSON(w, r, stats)
	}
}
