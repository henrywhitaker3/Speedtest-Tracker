package restore

import (
	"context"
	"net/http"
	"strings"
	"time"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/labstack/echo/v5"
)

type restoreRepo interface {
	Restore(context.Context, []results.Result) (int, int, error)
}

type restoreHandler struct {
	repo restoreRepo
}

type LegcayTime time.Time

func (ct *LegcayTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	*ct = LegcayTime(t)
	return nil
}

type LegacyResult struct {
	ID         int        `json:"id"`
	Ping       float64    `json:"ping"`
	Download   float64    `json:"download"`
	Upload     float64    `json:"upload"`
	CreatedAt  LegcayTime `json:"created_at"`
	UpdatedAt  LegcayTime `json:"updated_at"`
	ServerID   int        `json:"server_id"`
	ServerName string     `json:"server_name"`
	ServerHost string     `json:"server_host"`
	URL        string     `json:"url"`
	Failed     int        `json:"failed"`
}

type RestoreResponse struct {
	Restored int `json:"restored"`
	Skipped  int `json:"skipped"`
}

func (r *restoreHandler) Handler() handlers.Handler[[]LegacyResult, RestoreResponse] {
	return func(c *echo.Context, req []LegacyResult) (*RestoreResponse, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Restore results")
		defer span.End()

		var empty time.Time
		countEmpty := 0

		mapped := []results.Result{}
		for _, r := range req {
			if time.Time(r.CreatedAt).Equal(empty) {
				countEmpty++
				continue
			}
			id, _ := uuid.OrderedAt(time.Time(r.CreatedAt))
			mapped = append(mapped, results.Result{
				ID:         id,
				Ping:       &r.Ping,
				Download:   new(convertToBytes(r.Download)),
				Upload:     new(convertToBytes(r.Upload)),
				ServerID:   &r.ServerID,
				ServerHost: &r.ServerHost,
				ServerName: &r.ServerName,
				CreatedAt:  time.Time(r.CreatedAt),
				UpdatedAt:  time.Time(r.UpdatedAt),
				URL:        &r.URL,
				Complete:   true,
				Error:      failedToError(r.Failed),
			})
		}

		restored, skipped, err := r.repo.Restore(ctx, mapped)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &RestoreResponse{
			Restored: restored,
			Skipped:  skipped + countEmpty,
		}, nil
	}
}

func failedToError(failed int) *string {
	if failed == 0 {
		return nil
	}
	return new("unknown error")
}

func convertToBytes(input float64) int64 {
	return int64((input * 1000000) / 8)
}

func (r *restoreHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (r *restoreHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:           "Restore legacy backup",
		Tag:            "Restore",
		ID:             new("restoreJSON"),
		Description:    "Restore legacy speedtest json backups",
		Method:         http.MethodPost,
		Path:           "/v1/restore/json",
		Code:           http.StatusOK,
		GenerateSpec:   true,
		SkipValidation: true,
		Auths: []handlers.Auth{
			{
				Enabled: true,
				Name:    authapi.SessionSecurityScheme,
			},
		},
	}
}
