package compat

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/units"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/labstack/echo/v5"
)

type latestRepo interface {
	Latest(context.Context) (*results.Result, error)
}

type latestHandler struct {
	repo latestRepo
}

type LatestRequest struct{}

type LatestResponse struct {
	ID                 uuid.UUID                `json:"id"`
	Service            string                   `json:"service"`
	Ping               *float64                 `json:"ping"`
	Download           *int64                   `json:"download"`
	Upload             *int64                   `json:"upload"`
	DownloadBits       *int64                   `json:"download_bits,omitempty"`
	UploadBits         *int64                   `json:"upload_bits,omitempty"`
	DownloadBitsHuman  *string                  `json:"download_bits_human,omitempty"`
	UploadBitsHuman    *string                  `json:"upload_bits_human,omitempty"`
	DownloadBytes      *int                     `json:"download_bytes"`
	UploadBytes        *int                     `json:"upload_bytes"`
	DownloadBytesHuman *string                  `json:"download_bytes_human,omitempty"`
	UploadBytesHuman   *string                  `json:"upload_bytes_human,omitempty"`
	Benchmarks         any                      `json:"benchmarks"`
	Healthy            bool                     `json:"healthy"`
	Status             string                   `json:"status"`
	Scheduled          bool                     `json:"scheduled"`
	DispatchedBy       any                      `json:"dispatched_by"`
	Comments           any                      `json:"comments"`
	Data               *results.SpeedtestOutput `json:"data"`
	CreatedAt          string                   `json:"created_at"`
	UpdatedAt          string                   `json:"updated_at"`
}

// LatestResponseEnvelope matches Laravel's default JSON resource wrapper.
type LatestResponseEnvelope struct {
	Data *LatestResponse `json:"data"`
}

func (l *latestHandler) Handler() handlers.Handler[LatestRequest, LatestResponseEnvelope] {
	return func(c *echo.Context, req LatestRequest) (*LatestResponseEnvelope, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Get latest result")
		defer span.End()

		latest, err := l.repo.Latest(ctx)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &LatestResponseEnvelope{Data: latestResponse(c, latest)}, nil
	}
}

func latestResponse(c *echo.Context, result *results.Result) *LatestResponse {
	response := &LatestResponse{
		ID:        result.ID,
		Service:   "ookla",
		Ping:      result.Ping,
		Download:  result.Download,
		Upload:    result.Upload,
		Healthy:   result.Complete && result.Error == nil,
		Status:    resultStatus(result),
		Scheduled: true,
		CreatedAt: result.CreatedAt.Format(time.DateTime),
		UpdatedAt: result.UpdatedAt.Format(time.DateTime),
	}
	if result.Raw != nil {
		response.Data = new(result.Raw.Data())
	}

	if response.Data != nil && !auth.IsAuthenticated(c) {
		response.Data.Sanitise()
	}

	if result.Download != nil && *result.Download != 0 {
		bits := units.BytesToBits(*result.Download)
		response.DownloadBits = &bits
		human := units.BitsHuman(bits)
		response.DownloadBitsHuman = &human
	}
	if result.Upload != nil && *result.Upload != 0 {
		bits := units.BytesToBits(*result.Upload)
		response.UploadBits = &bits
		human := units.BitsHuman(bits)
		response.UploadBitsHuman = &human
	}
	if result.Raw != nil {
		data := result.Raw.Data()
		response.DownloadBytes = &data.Download.Bytes
		response.UploadBytes = &data.Upload.Bytes
		if data.Download.Bytes != 0 {
			human := formatFileSize(data.Download.Bytes)
			response.DownloadBytesHuman = &human
		}
		if data.Upload.Bytes != 0 {
			human := formatFileSize(data.Upload.Bytes)
			response.UploadBytesHuman = &human
		}
	}

	return response
}

func resultStatus(result *results.Result) string {
	if result.Error != nil {
		return "failed"
	}
	if result.Complete {
		return "completed"
	}
	return "running"
}

func formatFileSize(bytes int) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	quantity := float64(bytes)
	unit := 0
	for math.Abs(quantity)/1024 > 0.9 && unit < len(units)-1 {
		quantity /= 1024
		unit++
	}

	return fmt.Sprintf("%s %s", formatInteger(quantity), units[unit])
}

func formatInteger(value float64) string {
	number := strconv.FormatInt(int64(math.Round(value)), 10)
	start := 0
	if number[0] == '-' {
		start = 1
	}
	for i := len(number) - 3; i > start; i -= 3 {
		number = number[:i] + "," + number[i:]
	}

	return number
}

func (l *latestHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *latestHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Latest results",
		Tag:          "Compatibility",
		Description:  "Compatible endpoint with alexjustesen/speedtest-tracker",
		ID:           new("compatLatest"),
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		Path:         "/api/v1/results/latest",
		GenerateSpec: true,
	}
}
