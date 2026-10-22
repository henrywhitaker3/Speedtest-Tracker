package speedtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func testContext(t *testing.T) *echo.Context {
	t.Helper()
	return echo.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
}
