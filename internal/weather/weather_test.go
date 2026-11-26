package weather_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryStore(t *testing.T) {
	tests := []struct {
		name         string
		createResult bool
		missingID    bool
		wantError    string
	}{
		{name: "stores weather for an existing result", createResult: true},
		{name: "rejects an empty result ID", missingID: true, wantError: "empty result id"},
		{name: "rejects a nonexistent result", wantError: "FOREIGN KEY constraint failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.Sqlite(
				filepath.Join(t.TempDir(), "weather.db") + "?_pragma=foreign_keys(1)",
			)
			require.NoError(t, err)
			require.NoError(t, results.Register(db))

			resultID := uuid.MustOrdered()
			if tt.createResult {
				require.NoError(t, db.Create(&results.Result{ID: resultID}).Error)
			}
			if tt.missingID {
				resultID = uuid.UUID{}
			}
			data := weather.Weather{
				ResultID:      resultID,
				Code:          61,
				Temperature:   12.5,
				Precipitation: 1.2,
				WindSpeed:     7.5,
			}

			err = weather.NewRepo(db).Store(context.Background(), &data)
			if tt.wantError != "" {
				require.ErrorContains(t, err, tt.wantError)
				var count int64
				require.NoError(t, db.Model(&weather.Weather{}).Count(&count).Error)
				assert.Zero(t, count)
				return
			}
			require.NoError(t, err)
			var stored weather.Weather
			require.NoError(t, db.First(&stored, "result_id = ?", resultID).Error)
			assert.Equal(t, data, stored)
		})
	}
}
