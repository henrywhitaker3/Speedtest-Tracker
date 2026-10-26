package speedtest

import (
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Result struct {
	ID         uuid.UUID                `json:"id"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	DeletedAt  gorm.DeletedAt           `json:"deleted_at"`
	Ping       *float64                 `json:"ping"`
	Download   *int64                   `json:"download"`
	Upload     *int64                   `json:"upload"`
	ServerID   *int                     `json:"server_id"`
	ServerHost *string                  `json:"server_host"`
	ServerName *string                  `json:"server_name"`
	URL        *string                  `json:"url"`
	Error      *string                  `json:"error,omitempty"`
	Complete   bool                     `json:"complete"`
	Raw        *results.SpeedtestOutput `json:"raw,omitempty"`
}

func convert(c *echo.Context, r results.Result) Result {
	out := Result{
		ID:         r.ID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
		Ping:       r.Ping,
		Download:   r.Download,
		Upload:     r.Upload,
		ServerID:   r.ServerID,
		ServerHost: r.ServerHost,
		ServerName: r.ServerName,
		URL:        r.URL,
		Error:      r.Error,
		Complete:   r.Complete,
	}

	var output *results.SpeedtestOutput
	if r.Raw != nil {
		output = new(r.Raw.Data())
	}

	if !auth.IsAuthenticated(c) && output != nil {
		out.URL = nil
		output.Sanitise()
	}

	out.Raw = output

	return out
}
