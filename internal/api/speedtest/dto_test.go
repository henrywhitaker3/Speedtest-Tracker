package speedtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/labstack/echo/v5"
	"gorm.io/datatypes"
)

func TestConvertIncludesRawOnlyForAdministratorsOrWhenAuthIsDisabled(t *testing.T) {
	raw := datatypes.NewJSONType(results.SpeedtestOutput{})
	result := results.Result{Raw: &raw}

	tests := []struct {
		name        string
		authEnabled bool
		admin       bool
		wantRaw     bool
	}{
		{name: "hides raw from anonymous users", authEnabled: true},
		{name: "hides raw from non-administrators", authEnabled: true, admin: false},
		{name: "includes raw for administrators", authEnabled: true, admin: true, wantRaw: true},
		{name: "includes raw when authentication is disabled", wantRaw: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := echo.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
			c.Set("auth", tt.authEnabled)
			if tt.admin {
				c.Set("admin", true)
			}

			got := convert(c, result)
			if (got.Raw != nil) != tt.wantRaw {
				t.Errorf("Raw included = %t, want %t", got.Raw != nil, tt.wantRaw)
			}
		})
	}
}
