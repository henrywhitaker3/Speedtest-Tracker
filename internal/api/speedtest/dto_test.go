package speedtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
		{name: "includes PII for authenticated users", authEnabled: true, authorization: "Bearer valid-token", wantPII: true},
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

			if rec.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
			}
			if got.Raw == nil {
				t.Fatal("raw output was omitted")
			}
			if got.Raw.Type != "result" || got.Raw.PacketLoss != 1.5 {
				t.Errorf("non-PII raw data = %+v, want type and packet loss preserved", got.Raw)
			}
			hasPII := got.Raw.Isp != nil || got.Raw.Interface != nil || got.Raw.Server != nil || got.Raw.Result != nil
			if hasPII != tt.wantPII {
				t.Errorf("raw output contains PII = %t, want %t", hasPII, tt.wantPII)
			}
		})
	}
}

type dtoTokenRepo struct{}

func (dtoTokenRepo) Validate(context.Context, string) error {
	return nil
}
