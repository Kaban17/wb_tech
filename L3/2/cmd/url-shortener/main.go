package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/http-server/handler/url"
	config "url-shortener/internal/parser"
	sqlite "url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	LOCAL       = "local"
	DEVELOPMENT = "dev"
	PRODUCTION  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Starting application", "env", cfg.Env)
	s, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init database", "error", err)
		os.Exit(1)
	}
	// TODO : init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/shorten", url.New(log, s))
	router.Get("/{alias}", url.Redirect(log, s))
	router.Get("/stats/daily", url.DailyStats(log, s))
	router.Get("/stats/weekly", url.WeeklyStats(log, s))
	router.Get("/stats/monthly", url.MonthlyStats(log, s))
	// TODO : run server
	serv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.IdleTimeout,
	}
	log.Info("Starting server", "address", cfg.HttpServer.Address)
	serv.ListenAndServe()

}
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case DEVELOPMENT:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case PRODUCTION:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
