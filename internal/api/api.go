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
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/compat"
	configapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/restore"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/speedtest"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/duration"
	"github.com/henrywhitaker3/windowframe/v2/http"
	mw "github.com/henrywhitaker3/windowframe/v2/http/middleware"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
	"gorm.io/gorm"
)

const serviceName = "Speedtest-Tracker"

func New(
	version string,
	conf *config.Config,
	db *gorm.DB,
	host *local.Host,
	runner *workflow.WorkflowService,
) (*http.HTTP, error) {
	openapiOpts := openapiOpts(version, conf.App.URL)
	mut := specMutations()
	if !conf.Auth.Enabled {
		openapiOpts.AuthConfig = nil
	} else {
		mut = append(
			mut,
			func(r *openapi3.Reflector) {
				r.SpecEns().SetAPIKeySecurity(
					authapi.SessionSecurityScheme,
					authapi.SessionCookieName,
					openapi.InCookie,
					"Browser session authentication.",
				)
			},
		)
	}
	h := http.New(http.HTTPOpts{
		Port:          8080,
		Openapi:       openapiOpts,
		Logger:        slog.With("component", "api"),
		SpecMutations: mut,
	})

	h.Use(middleware.Recover())
	h.Use(mw.Tracing(serviceName))
	h.Use(mw.Metrics(serviceName, prometheus.DefaultRegisterer))
	h.Use(mw.LoggerAt(slog.LevelDebug))
	h.Use(corsMiddleware(conf.App.UIURL))
	h.Use(staticMiddleware(conf.App.UIDir))
	for _, middleware := range authapi.Middleware(authapi.MiddlewareOpts{
		Enabled:     conf.Auth.Enabled,
		Service:     host.Service(),
		TokenRepo:   auth.NewRepo(db),
		AdminGroup:  conf.Auth.AdminGroup,
		Routes:      authRoutes(),
		AdminRoutes: adminRoutes(),
	}) {
		h.Use(middleware)
	}
	if err := registerRoutes(h, db, host, runner, conf); err != nil {
		return nil, err
	}
	registerValidators(h)

	return h, nil
}

func openapiOpts(version, publicURL string) http.OpenapiOpts {
	return http.OpenapiOpts{
		Enabled:        true,
		ServiceName:    serviceName,
		ServiceVersion: version,
		PublicURL:      publicURL,
		AuthConfig: []struct {
			Enabled     bool
			Name        string
			Format      string
			Description string
		}{
			{
				Enabled:     true,
				Name:        authapi.APITokenSecurityScheme,
				Format:      "API token",
				Description: "Send an API token as `Authorization: Bearer <token>`.",
			},
		},
	}
}

func specMutations() []http.SpecMutator {
	return []http.SpecMutator{
		func(r *openapi3.Reflector) {
			msgStatusDef := jsonschema.Schema{}
			msgStatusDef.AddType(jsonschema.String)
			statuses := []any{
				"pending",
				"running",
				"completed",
				"failed",
				"skipped",
				"compensating",
				"compensated",
				"compensation-failed",
			}
			msgStatusDef.WithEnum(statuses...)
			r.AddTypeMapping(workflow.StepStatus(""), msgStatusDef)
		},
		func(r *openapi3.Reflector) {
			msgStatusDef := jsonschema.Schema{}
			msgStatusDef.AddType(jsonschema.String)
			statuses := []any{
				config.PreRun,
				config.PostRun,
			}
			msgStatusDef.WithEnum(statuses...)
			r.AddTypeMapping(config.NotificationPhase(""), msgStatusDef)
		},
		func(r *openapi3.Reflector) {
			msgStatusDef := jsonschema.Schema{}
			msgStatusDef.AddType(jsonschema.String)
			statuses := []any{
				config.Healthcheck,
				config.Webhook,
			}
			msgStatusDef.WithEnum(statuses...)
			r.AddTypeMapping(config.NotificationKind(""), msgStatusDef)
		},
	}
}

func authRoutes() authapi.ProtectedRoutes {
	return authapi.ProtectedRoutes{
		{Method: stdhttp.MethodPost, Path: "/v1/speedtest/run"},
		{Method: stdhttp.MethodDelete, Path: "/v1/speedtest/purge"},
		{Method: stdhttp.MethodDelete, Path: "/v1/speedtest/:id"},
		{Method: stdhttp.MethodGet, Path: "/v1/auth/tokens"},
		{Method: stdhttp.MethodPost, Path: "/v1/auth/tokens"},
		{Method: stdhttp.MethodDelete, Path: "/v1/auth/tokens/:id"},
		{Method: stdhttp.MethodGet, Path: "/v1/config/notifications"},
		{Method: stdhttp.MethodPost, Path: "/v1/restore/json"},
	}
}

func adminRoutes() authapi.ProtectedRoutes {
	return authapi.ProtectedRoutes{
		{Method: stdhttp.MethodDelete, Path: "/v1/speedtest/purge"},
		{Method: stdhttp.MethodDelete, Path: "/v1/speedtest/:id"},
		{Method: stdhttp.MethodGet, Path: "/v1/auth/tokens"},
		{Method: stdhttp.MethodPost, Path: "/v1/auth/tokens"},
		{Method: stdhttp.MethodDelete, Path: "/v1/auth/tokens/:id"},
		{Method: stdhttp.MethodGet, Path: "/v1/config/notifications"},
		{Method: stdhttp.MethodPost, Path: "/v1/restore/json"},
	}
}

func staticMiddleware(root string) echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  root,
		HTML5: true,
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path
			return path == "/v1" || strings.HasPrefix(path, "/v1/") ||
				strings.HasPrefix(path, "/docs") || strings.HasPrefix(path, "/api")
		},
	})
}

func registerRoutes(
	h *http.HTTP,
	db *gorm.DB,
	host *local.Host,
	runner *workflow.WorkflowService,
	conf *config.Config,
) error {
	if err := authapi.RegisterRoutes(h, host, db, conf); err != nil {
		return fmt.Errorf("register auth routes: %w", err)
	}
	if err := speedtest.Register(h, db, host, runner); err != nil {
		return err
	}
	configapi.RegisterRoutes(h, conf)
	restore.Register(h, db)
	compat.Register(h, db)
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
