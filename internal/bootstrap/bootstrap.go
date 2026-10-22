// Package bootstrap
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
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
) (*gorm.DB, *local.Host, error) {
	_, err := tracing.InitTracer(tracing.TracingOpts{
		ServiceName: "Speedtest-Tracker",
	}, version)
	if err != nil {
		return nil, nil, fmt.Errorf("init tracer: %w", err)
	}

	slog.Debug("loaded config", "config", conf)

	if err := ookla.EnsureBinary(ctx, ooklaVersion); err != nil {
		return nil, nil, fmt.Errorf("download speedtest-cli: %w", err)
	}

	db, err := database.Open(conf.Database.Driver, string(conf.Database.URL))
	if err != nil {
		return nil, nil, fmt.Errorf("connect to db: %w", err)
	}
	if err := results.Register(db); err != nil {
		return nil, nil, fmt.Errorf("register results model: %w", err)
	}

	host, err := runtime(conf.Database, conf.Runtime)
	if err != nil {
		return nil, nil, fmt.Errorf("init runtime: %w", err)
	}

	if err := ookla.RegisterActor(host, results.NewRepo(db)); err != nil {
		return nil, nil, fmt.Errorf("register ookla actors: %w", err)
	}

	return db, host, nil
}

func runtime(db config.Database, conf config.Runtime) (*local.Host, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
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
