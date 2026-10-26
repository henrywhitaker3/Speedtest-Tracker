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

func TestConvertIncludesRawOnlyForAuthenticatedUsersOrWhenAuthIsDisabled(t *testing.T) {
	raw := datatypes.NewJSONType(results.SpeedtestOutput{})
	result := results.Result{Raw: &raw}

	tests := []struct {
		name          string
		authEnabled   bool
		authorization string
		wantRaw       bool
	}{
		{name: "hides raw from anonymous users", authEnabled: true},
		{name: "includes raw for authenticated users", authEnabled: true, authorization: "Bearer valid-token", wantRaw: true},
		{name: "includes raw when authentication is disabled", wantRaw: true},
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
			if (got.Raw != nil) != tt.wantRaw {
				t.Errorf("Raw included = %t, want %t", got.Raw != nil, tt.wantRaw)
			}
		})
	}
}

type dtoTokenRepo struct{}

func (dtoTokenRepo) Validate(context.Context, string) error {
	return nil
}
