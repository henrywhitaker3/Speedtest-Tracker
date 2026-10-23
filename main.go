package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/bootstrap"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/meta"
)

const (
	version      = "dev"
	ooklaVersion = "1.2.0"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	db, host, err := bootstrap.Bootstrap(ctx, conf, version, ooklaVersion)
	if err != nil {
		slog.Error("bootstrap app", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := meta.New().Run(ctx); err != nil {
			slog.Error("run meta server", "error", err)
			os.Exit(1)
		}
	}()

	http, err := api.New(
		version,
		conf,
		db,
		host,
	)
	if err != nil {
		slog.Error("init api", "error", err)
		os.Exit(1)
	}

	go func() {
		<-host.Ready()
		meta.Ready = true
		if err := http.Start(ctx); err != nil {
			slog.Error("run api", "error", err)
			os.Exit(1)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		_ = http.Stop(ctx)
	}()

	if err := host.Run(ctx); err != nil {
		slog.Error("start runtime", "error", err)
		os.Exit(1)
	}
}
