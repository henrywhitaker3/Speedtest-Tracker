package config

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListServersHandlerListsServers(t *testing.T) {
	servers := []ookla.Server{{ID: 1, Name: "Speedtest"}}
	lister := &serverListerStub{servers: servers}

	got, err := (&listServersHandler{lister: lister}).Handler()(testContext(t), nil)

	require.NoError(t, err)
	assert.Equal(t, servers, *got)
	assert.True(t, lister.called)
}

func TestListServersHandlerWrapsListerError(t *testing.T) {
	want := errors.New("list servers")

	_, err := (&listServersHandler{lister: &serverListerStub{err: want}}).Handler()(
		testContext(t),
		nil,
	)

	assert.ErrorIs(t, err, want)
}

type serverListerStub struct {
	servers []ookla.Server
	err     error
	called  bool
}

func (s *serverListerStub) ListServers(context.Context) ([]ookla.Server, error) {
	s.called = true
	return s.servers, s.err
}

func testContext(t *testing.T) *echo.Context {
	t.Helper()
	return echo.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
}
