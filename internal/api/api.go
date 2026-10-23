// Package api
package api

import (
	"fmt"
	"log/slog"
	stdhttp "net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	configapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/speedtest"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/duration"
	"github.com/henrywhitaker3/windowframe/v2/http"
	mw "github.com/henrywhitaker3/windowframe/v2/http/middleware"
	"github.com/italypaleale/francis/host/local"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func New(
	version string,
	conf *config.Config,
	db *gorm.DB,
	host *local.Host,
) (*http.HTTP, error) {
	h := http.New(http.HTTPOpts{
		Port: 8080,
		Openapi: http.OpenapiOpts{
			Enabled:        true,
			ServiceName:    "Speedtest-Tracker",
			ServiceVersion: version,
			PublicURL:      conf.App.URL,
		},
		Logger: slog.With("component", "api"),
	})

	h.Use(middleware.Recover())
	h.Use(mw.Tracing("Speedtest-Tracker"))
	h.Use(mw.Metrics("Speedtest-Tracker", prometheus.DefaultRegisterer))
	h.Use(mw.Logger())
	h.Use(corsMiddleware(conf.App.UIURL))
	h.Use(staticMiddleware(conf.App.UIDir))

	if conf.Auth.Enabled {
		h.Use(authapi.SessionMiddleware())
		h.Use(authapi.AuthMiddleware(authapi.AuthMiddlewareOpts{
			Service:   host.Service(),
			TokenRepo: auth.NewRepo(db),
			Routes: []string{
				"/v1/speedtest/run",
				"/v1/auth/tokens",
			},
		}))
	}

	if err := registerRoutes(h, db, host, conf); err != nil {
		return nil, err
	}
	registerValidators(h)

	return h, nil
}

func staticMiddleware(root string) echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  root,
		HTML5: true,
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path
			return path == "/v1" || strings.HasPrefix(path, "/v1/") ||
				strings.HasPrefix(path, "/docs")
		},
	})
}

func registerRoutes(
	h *http.HTTP,
	db *gorm.DB,
	host *local.Host,
	conf *config.Config,
) error {
	if err := authapi.RegisterRoutes(h, host, db, conf); err != nil {
		return fmt.Errorf("register auth routes: %w", err)
	}
	if err := speedtest.Register(h, db, host); err != nil {
		return err
	}
	configapi.RegisterRoutes(h, conf)
	return nil
}

func registerValidators(h *http.HTTP) {
	validTimeseriesPeriods := []time.Duration{
		0, // Allows for empty default
		time.Hour * 24,
		time.Hour * 24 * 7,
		time.Hour * 24 * 30,
	}
	h.Validator.RegisterValidation("timeseries_period", func(fl validator.FieldLevel) bool {
		dur := fl.Field().Interface().(duration.StringDuration)
		return slices.Contains(validTimeseriesPeriods, dur.Duration())
	})
}

func corsMiddleware(uiURL string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{uiURL},
		AllowCredentials: true,
		AllowMethods: []string{
			stdhttp.MethodGet,
			stdhttp.MethodPost,
			stdhttp.MethodOptions,
			stdhttp.MethodDelete,
		},
		AllowHeaders: []string{echo.HeaderContentType},
		MaxAge:       86_400,
	})
}
