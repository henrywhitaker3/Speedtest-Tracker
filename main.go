package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/api"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/bootstrap"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/config"
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

	host, err := bootstrap.Bootstrap(ctx, conf, version, ooklaVersion)
	if err != nil {
		slog.Error("bootstrap app", "error", err)
		os.Exit(1)
	}

	api := api.New(conf.App.URL, version, host.Service())
	go func() {
		if err := api.Start(ctx); err != nil {
			slog.Error("run api", "error", err)
			os.Exit(1)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		_ = api.Stop(ctx)
	}()

	if err := host.Run(ctx); err != nil {
		slog.Error("start runtime", "error", err)
		os.Exit(1)
	}
}
