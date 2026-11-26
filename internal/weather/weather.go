// Package weather
package weather

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

type WeatherData struct {
	Code          WeatherCode `json:"code"`
	Temperature   float64     `json:"temperature"`
	Precipitation float64     `json:"precipitation"`
	WindSpeed     float64     `json:"wind_speed"`
}

type Weather struct {
	ResultID uuid.UUID `gorm:"type:uuid;primaryKey" json:"result_id"`
	WeatherData
}

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Store(ctx context.Context, data *Weather) error {
	var empty uuid.UUID
	if data.ResultID == empty {
		return fmt.Errorf("empty result id")
	}
	if err := gorm.G[Weather](r.db).Create(ctx, data); err != nil {
		return fmt.Errorf("store weather data: %w", err)
	}
	return nil
}
