package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OSShip/mentors/internal/config"
	"github.com/OSShip/mentors/internal/events"
	"github.com/OSShip/mentors/internal/github"
	"github.com/OSShip/mentors/internal/handler"
	"github.com/OSShip/mentors/internal/store"
	"github.com/OSShip/utils/observability"
)

func main() {
	cfg := config.Load()
	observability.InitSentry("mentors")
	defer observability.FlushSentry(2 * time.Second)
	logger := observability.InitLogger("mentors")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connected")

	pub := events.New(cfg.KafkaBrokers)
	defer pub.Close()
	logger.Info("kafka publisher ready", "brokers", cfg.KafkaBrokers)

	h := &handler.Handler{
		Store:  store.New(pool),
		Events: pub,
		Github: github.New(cfg.GithubToken),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(observability.SentryHTTPMiddleware)
	r.Use(observability.SentryRecoverMiddleware("mentors"))
	r.Use(observability.SentryErrorMiddleware("mentors"))
	r.Use(observability.RequestLogMiddleware("mentors"))
	r.Use(observability.PrometheusMiddleware("mentors"))

	r.Get("/health", observability.HealthHandler("mentors"))
	r.Get("/metrics", observability.MetricsHandler().ServeHTTP)

	r.Post("/apply", h.Apply)
	r.Get("/admin/applications", h.ListApplications)
	r.Patch("/admin/applications/{id}", h.ReviewApplication)

	logger.Info("mentors listening", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}
