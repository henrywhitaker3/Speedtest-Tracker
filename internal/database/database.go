// Package database
package database

import (
	"fmt"
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func Open(driver string, url string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch driver {
	case "sqlite":
		db, err = Sqlite(url)
	case "postgres":
		db, err = Postgres(url)
	default:
		return nil, fmt.Errorf("unsupported db driver: %s", driver)
	}
	if err != nil {
		return nil, fmt.Errorf("open gorm db: %w", err)
	}

	if err := instrument(db); err != nil {
		return nil, fmt.Errorf("instrument db: %w", err)
	}

	return db, nil
}

func Sqlite(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), config())
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}
	return db, nil
}

func Postgres(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url))
	if err != nil {
		return nil, fmt.Errorf("open postgres db: %w", err)
	}
	return db, nil
}

func instrument(db *gorm.DB) error {
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return fmt.Errorf("register sqlite tracer: %w", err)
	}
	return nil
}

func config() *gorm.Config {
	return &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(slog.With("component", "db").Handler()),
			slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, slog.LevelWarn),
		),
	}
}
