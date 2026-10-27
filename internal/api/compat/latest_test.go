package compat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/labstack/echo/v5"
	"gorm.io/datatypes"
)

func TestLatestResponse(t *testing.T) {
	download := int64(1_000_000)
	upload := int64(500)
	ping := 12.5
	createdAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Minute)
	data := results.SpeedtestOutput{}
	data.Download.Bytes = 2_048
	data.Upload.Bytes = 512
	raw := datatypes.NewJSONType(data)

	context := echo.NewContext(
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRecorder(),
	)
	response := latestResponse(context, &results.Result{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Ping:      &ping,
		Download:  &download,
		Upload:    &upload,
		Complete:  true,
		Raw:       &raw,
	})

	assert.Equal(t, "ookla", response.Service)
	assert.Equal(t, "completed", response.Status)
	assert.True(t, response.Healthy)
	assert.True(t, response.Scheduled)
	assert.Equal(t, int64(8_000_000), *response.DownloadBits)
	assert.Equal(t, "8.00 Mbps", *response.DownloadBitsHuman)
	assert.Equal(t, "4.00 kbps", *response.UploadBitsHuman)
	assert.Equal(t, "2 KB", *response.DownloadBytesHuman)
	assert.Equal(t, "512 B", *response.UploadBytesHuman)
	assert.Equal(t, "2026-01-02 03:04:05", response.CreatedAt)
	assert.Equal(t, "2026-01-02 03:05:05", response.UpdatedAt)
}

func TestLatestResponseWithoutRawData(t *testing.T) {
	context := echo.NewContext(
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRecorder(),
	)

	response := latestResponse(context, &results.Result{})
	assert.Nilf(t, response.Data, "data = %#v, want nil", response.Data)
}

func TestLatestResponseEnvelope(t *testing.T) {
	download := int64(1_000_000)
	response := LatestResponseEnvelope{Data: &LatestResponse{Download: &download}}

	encoded, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded struct {
		Data struct {
			Download int64 `json:"download"`
		} `json:"data"`
	}
	require.Nil(t, json.Unmarshal(encoded, &decoded))
	assert.Equalf(
		t,
		download,
		decoded.Data.Download,
		"data.download = %d, want %d",
		decoded.Data.Download,
		download,
	)
}

func TestResultStatus(t *testing.T) {
	tests := []struct {
		name   string
		result results.Result
		want   string
	}{
		{name: "running", want: "running"},
		{name: "completed", result: results.Result{Complete: true}, want: "completed"},
		{name: "failed", result: results.Result{Error: new("failed")}, want: "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, resultStatus(&tt.result))
		})
	}
}
