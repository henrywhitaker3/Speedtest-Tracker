// Package results
package results

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

type Result struct {
	database.BaseModel

	Ping       *float64 `json:"ping"`
	Download   *float64 `json:"download"`
	Upload     *float64 `json:"upload"`
	ServerID   *int     `json:"server_id"`
	ServerHost *string  `json:"server_host"`
	ServerName *string  `json:"server_name"`
	URL        *string  `json:"url"`
	Error      *string  `json:"error,omitempty"`
	Complete   bool     `json:"complete"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context) (*Result, error) {
	result := &Result{
		ID: uuid.MustOrdered(),
	}
	if err := gorm.G[Result](r.db).Create(ctx, result); err != nil {
		return nil, fmt.Errorf("store initial result: %w", err)
	}

	return result, nil
}

func (r *Repository) Fresh(ctx context.Context, result *Result) error {
	fresh, err := r.Get(ctx, result.ID)
	if err != nil {
		return err
	}
	*result = *fresh
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Result, error) {
	result, err := gorm.G[Result](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get result: %w", err)
	}
	return &result, nil
}

func (r *Repository) Update(ctx context.Context, result *Result) error {
	if _, err := gorm.G[Result](r.db).
		Where("id = ?", result.ID).
		Select("*").
		Updates(ctx, *result); err != nil {
		return fmt.Errorf("store updated result: %w", err)
	}
	return nil
}

func Register(db *gorm.DB) error {
	return db.AutoMigrate(
		&Result{},
	)
}
