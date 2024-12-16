package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mihailtudos/url-shortener/internal/config"
	"github.com/mihailtudos/url-shortener/internal/lib/logger/sl"
	"github.com/mihailtudos/url-shortener/internal/storage/sqlite"

	"github.com/mihailtudos/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/mihailtudos/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/mihailtudos/url-shortener/internal/http-server/middleware/logger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()

	log.Printf("config: %+v", cfg)

	// init logger
	logger := setupLogger(cfg.Env)

	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	// init router
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

    router.Route("/url", func(r chi.Router) {
        r.Use(middleware.BasicAuth("url-shortener", map[string]string{
            cfg.HTTPServer.User: cfg.HTTPServer.Password,
        }))
        
        r.Post("/", save.New(logger, storage))
    })
	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
}

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
