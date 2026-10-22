// Package api
package api

import (
	"log/slog"
	stdhttp "net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/api/speedtest"
	"github.com/henrywhitaker3/windowframe/v2/http"
	mw "github.com/henrywhitaker3/windowframe/v2/http/middleware"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func New(url, version string, db *gorm.DB, service *actor.Service) *http.HTTP {
	h := http.New(http.HTTPOpts{
		Port: 8080,
		Openapi: http.OpenapiOpts{
			Enabled:        true,
			ServiceName:    "Speedtest-Tracker",
			ServiceVersion: version,
			PublicURL:      url,
		},
		Logger: slog.With("component", "api"),
	})

	h.Use(mw.Tracing("Speedtest-Tracker"))
	h.Use(mw.Metrics("Speedtest-Tracker", prometheus.DefaultRegisterer))
	h.Use(mw.Logger())
	h.Use(corsMiddleware())

	registerRoutes(h, db, service)

	return h
}

func registerRoutes(h *http.HTTP, db *gorm.DB, service *actor.Service) {
	speedtest.Register(h, db, service)
}

func corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			stdhttp.MethodGet,
			stdhttp.MethodPost,
			stdhttp.MethodOptions,
		},
		AllowHeaders: []string{echo.HeaderContentType},
		MaxAge:       86_400,
	})
}
