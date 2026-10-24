// Package bootstrap
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/italypaleale/francis/components/postgres"
	"github.com/italypaleale/francis/components/sqlite"
	"github.com/italypaleale/francis/host/local"
	"gorm.io/gorm"
)

func Bootstrap(
	ctx context.Context,
	conf *config.Config,
	version, ooklaVersion string,
) (*http.HTTP, *gorm.DB, *local.Host, error) {
	_, err := tracing.InitTracer(tracing.TracingOpts{
		ServiceName: "Speedtest-Tracker",
		SampleRate:  1,
	}, version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("init tracer: %w", err)
	}

	slog.Debug("loaded config", "config", conf)

	if err := ookla.EnsureBinary(ctx, ooklaVersion); err != nil {
		return nil, nil, nil, fmt.Errorf("download speedtest-cli: %w", err)
	}

	db, err := database.Open(conf.Database.Driver, string(conf.Database.URL))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connect to db: %w", err)
	}
	if err := results.Register(db); err != nil {
		return nil, nil, nil, fmt.Errorf("register results model: %w", err)
	}
	if err := auth.Register(db); err != nil {
		return nil, nil, nil, fmt.Errorf("register auth models: %w", err)
	}

	host, err := runtime(conf.Database, conf.Runtime)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("init runtime: %w", err)
	}

	resultsRepo := results.NewRepo(db)

	workflow, err := ookla.RegisterActor(ookla.RegisterActorOpts{
		Host:         host,
		ServerID:     conf.Server,
		Runner:       ookla.Run,
		WorkflowRepo: resultsRepo,
		CronjobRepo:  resultsRepo,
		Schedule:     conf.Schedule,
		Healthcheck:  conf.Notifications.Healthcheck,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("register ookla actors: %w", err)
	}
	if err := results.RegisterActors(host, resultsRepo, conf.RetentionDays.Duration()); err != nil {
		return nil, nil, nil, fmt.Errorf("register results actors: %w", err)
	}
	if conf.Auth.Enabled {
		pctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		provider, err := oidc.NewProvider(pctx, conf.Auth.IssuerURL)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("build oidc provider: %w", err)
		}
		if err := auth.RegisterActor(
			conf.Auth,
			auth.CallbackURL(conf.App.URL),
			host,
			provider,
		); err != nil {
			return nil, nil, nil, fmt.Errorf("register auth actor: %w", err)
		}
	}

	http, err := api.New(
		version,
		conf,
		db,
		host,
		workflow,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create api handler: %w", err)
	}

	return http, db, host, nil
}

func runtime(db config.Database, conf config.Runtime) (*local.Host, error) {
	hostname, err := hostname(conf)
	if err != nil {
		return nil, err
	}

	var dbOpt local.HostOption
	switch db.Driver {
	case "sqlite":
		dbOpt = local.WithSQLiteProvider(sqlite.SQLiteProviderOptions{
			ConnectionString: string(db.URL),
		})
	case "postgres":
		dbOpt = local.WithPostgresProvider(postgres.PostgresProviderOptions{
			ConnectionString: string(db.URL),
		})
	}

	host, err := local.NewHost(
		local.WithAddress(net.JoinHostPort(hostname, "7157")),
		local.WithRuntimePSKs([]byte(conf.PSK)),
		local.WithLogger(slog.With("component", "runtime")),
		dbOpt,
	)
	if err != nil {
		return nil, fmt.Errorf("init host: %w", err)
	}
	return host, nil
}

func hostname(conf config.Runtime) (string, error) {
	if conf.Hostname != "" {
		return conf.Hostname, nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("get hostname: %w", err)
	}
	return hostname, nil
}
