package compat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
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
		BaseModel: database.BaseModel{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Ping:     &ping,
		Download: &download,
		Upload:   &upload,
		Complete: true,
		Raw:      &raw,
	})

	if response.Service != "ookla" || response.Status != "completed" || !response.Healthy || !response.Scheduled {
		t.Fatalf("unexpected result metadata: %+v", response)
	}
	if got := *response.DownloadBits; got != 8_000_000 {
		t.Errorf("download_bits = %d, want 8000000", got)
	}
	if got := *response.DownloadBitsHuman; got != "8.00 Mbps" {
		t.Errorf("download_bits_human = %q, want %q", got, "8.00 Mbps")
	}
	if got := *response.UploadBitsHuman; got != "4.00 kbps" {
		t.Errorf("upload_bits_human = %q, want %q", got, "4.00 kbps")
	}
	if got := *response.DownloadBytesHuman; got != "2 KB" {
		t.Errorf("download_bytes_human = %q, want %q", got, "2 KB")
	}
	if got := *response.UploadBytesHuman; got != "512 B" {
		t.Errorf("upload_bytes_human = %q, want %q", got, "512 B")
	}
	if response.CreatedAt != "2026-01-02 03:04:05" || response.UpdatedAt != "2026-01-02 03:05:05" {
		t.Errorf("timestamps = %q, %q", response.CreatedAt, response.UpdatedAt)
	}
}

func TestLatestResponseWithoutRawData(t *testing.T) {
	context := echo.NewContext(
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRecorder(),
	)

	response := latestResponse(context, &results.Result{})
	if response.Data != nil {
		t.Errorf("data = %#v, want nil", response.Data)
	}
}

func TestLatestResponseEnvelope(t *testing.T) {
	download := int64(1_000_000)
	response := LatestResponseEnvelope{Data: &LatestResponse{Download: &download}}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		Data struct {
			Download int64 `json:"download"`
		} `json:"data"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Data.Download != download {
		t.Errorf("data.download = %d, want %d", decoded.Data.Download, download)
	}
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
			if got := resultStatus(&tt.result); got != tt.want {
				t.Errorf("resultStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}
