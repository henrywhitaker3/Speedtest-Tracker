// Package weather
package weather

import (
	"context"
	"fmt"
	"time"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WeatherData struct {
	// WMO weather interpretation code.
	Code WeatherCode `json:"code"`
	// Degrees Celsius.
	Temperature float64 `json:"temperature"`
	// Millimetres of rain and snow over the past hour.
	Precipitation float64 `json:"precipitation"`
	// Metres per second.
	WindSpeed float64 `json:"wind_speed"`
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

func (r *Repository) WithoutWeather(ctx context.Context) (map[uuid.UUID]time.Time, error) {
	var rows []struct {
		ID        uuid.UUID
		CreatedAt time.Time
	}
	err := r.db.WithContext(ctx).Table("results").
		Select("results.id, results.created_at").
		Joins("LEFT JOIN weathers ON weathers.result_id = results.id").
		Where("results.deleted_at IS NULL AND weathers.result_id IS NULL").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find results without weather: %w", err)
	}
	out := make(map[uuid.UUID]time.Time, len(rows))
	for _, row := range rows {
		out[row.ID] = row.CreatedAt
	}
	return out, nil
}

func (r *Repository) Store(ctx context.Context, data []Weather) error {
	for _, d := range data {
		var empty uuid.UUID
		if d.ResultID == empty {
			return fmt.Errorf("empty result id")
		}
	}

	if err := gorm.G[Weather](
		r.db,
		clause.OnConflict{DoNothing: true},
	).CreateInBatches(ctx, &data, 100); err != nil {
		return fmt.Errorf("store weather data: %w", err)
	}
	return nil
}

func RegisterDB(db *gorm.DB) error {
	return db.AutoMigrate(&Weather{})
}
