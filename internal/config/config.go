// Package config
package config

import (
	"fmt"
	"slices"

	"github.com/adhocore/gronx"
	"github.com/henrywhitaker3/windowframe/v2/config"
	"github.com/henrywhitaker3/windowframe/v2/log"
	"github.com/joho/godotenv"
)

type App struct {
	URL   string `env:"URL,default=http://localhost:8080"`
	UIDir string `env:"UI_DIR,default=ui"`
}

type Database struct {
	Driver string                `env:"DRIVER, default=sqlite"`
	URL    config.RedactedString `env:"URL, default=db.sqlite"`
}

type Runtime struct {
	Hostname string                `env:"HOSTNAME"`
	PSK      config.RedactedString `env:"PSK"`
}

type Config struct {
	App App `env:",prefix=APP_"`

	LogLevel   string `env:"LOG_LEVEL"`
	AcceptUELA bool   `env:"ACCEPT_EULA"`

	Database Database `env:",prefix=DB_"`
	Runtime  Runtime  `env:",prefix=RUNTIME_"`

	Schedule string `env:"TEST_SCHEDULE, default=5 */2 * * *"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	conf, err := config.NewParser[Config]().WithExtractors(
		config.NewEnvExtractor[Config](),
	).Parse()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	log.Setup(log.SlogLevel(conf.LogLevel))

	if !conf.AcceptUELA {
		return nil, fmt.Errorf("must accep eula")
	}
	if !slices.Contains([]string{"sqlite", "postgres"}, conf.Database.Driver) {
		return nil, fmt.Errorf("unknown db driver %s", conf.Database.Driver)
	}
	if !gronx.IsValid(conf.Schedule) {
		return nil, fmt.Errorf("invalid cron schedule %s", conf.Schedule)
	}

	return &conf, nil
}
