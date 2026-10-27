package speedtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/labstack/echo/v5"
	"gorm.io/datatypes"
)

func TestConvertSanitisesRawOutputForUnauthenticatedUsers(t *testing.T) {
	isp := "Example ISP"
	raw := datatypes.NewJSONType(results.SpeedtestOutput{
		Type:       "result",
		PacketLoss: 1.5,
		Isp:        &isp,
		Interface: &results.Interface{
			InternalIP: "192.168.1.2",
			ExternalIP: "203.0.113.1",
		},
		Server: &results.Server{Host: "speedtest.example.com"},
		Result: &results.OoklaResult{URL: "https://www.speedtest.net/result/1"},
	})
	result := results.Result{Raw: &raw}

	tests := []struct {
		name          string
		authEnabled   bool
		authorization string
		wantPII       bool
	}{
		{name: "removes PII from anonymous users", authEnabled: true},
		{
			name:          "includes PII for authenticated users",
			authEnabled:   true,
			authorization: "Bearer valid-token",
			wantPII:       true,
		},
		{name: "includes PII when authentication is disabled", wantPII: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(authapi.Middleware(authapi.MiddlewareOpts{
				Enabled:   tt.authEnabled,
				TokenRepo: dtoTokenRepo{},
			})...)

			var got Result
			e.GET("/v1/speedtest/latest", func(c *echo.Context) error {
				got = convert(c, result)
				return c.NoContent(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/v1/speedtest/latest", nil)
			if tt.authorization != "" {
				req.Header.Set(echo.HeaderAuthorization, tt.authorization)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equalf(
				t,
				http.StatusNoContent,
				rec.Code,
				"status = %d, want %d",
				rec.Code,
				http.StatusNoContent,
			)
			require.NotNilf(t, got.Raw, "raw output was omitted")
			assert.Equal(t, "result", got.Raw.Type)
			assert.Equal(t, 1.5, got.Raw.PacketLoss)
			if !tt.wantPII {
				assert.Nil(t, got.Raw.Isp)
				assert.Nil(t, got.Raw.Interface)
				assert.Nil(t, got.Raw.Server)
				assert.Nil(t, got.URL)
				assert.Nil(t, got.ServerHost)
				assert.Nil(t, got.ServerID)
				assert.Nil(t, got.ServerName)
			}
		})
	}
}

type dtoTokenRepo struct{}

func (dtoTokenRepo) Validate(context.Context, string) error {
	return nil
}
