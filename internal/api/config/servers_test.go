package config

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListServersHandlerListsServers(t *testing.T) {
	servers := []ookla.Server{{ID: 1, Name: "Speedtest"}}
	invoker := &serverListerInvokerStub{
		envelope: serverListEnvelope{servers: servers},
	}

	got, err := (&listServersHandler{invoker: invoker}).Handler()(testContext(t), nil)

	require.NoError(t, err)
	assert.Equal(t, servers, *got)
	assert.True(t, invoker.called)
	assert.Equal(t, ookla.Actor, invoker.actorType)
	assert.Equal(t, actor.SingletonActorID, invoker.actorID)
	assert.Equal(t, ookla.ListServers, invoker.method)
	assert.Nil(t, invoker.data)
}

func TestListServersHandlerWrapsInvokeError(t *testing.T) {
	want := errors.New("invoke list servers")

	_, err := (&listServersHandler{invoker: &serverListerInvokerStub{err: want}}).Handler()(
		testContext(t),
		nil,
	)

	assert.ErrorIs(t, err, want)
}

func TestListServersHandlerWrapsDecodeError(t *testing.T) {
	want := errors.New("decode servers")

	_, err := (&listServersHandler{invoker: &serverListerInvokerStub{
		envelope: serverListEnvelope{err: want},
	}}).Handler()(testContext(t), nil)

	assert.ErrorIs(t, err, want)
}

type serverListerInvokerStub struct {
	envelope  actor.Envelope
	err       error
	actorType string
	actorID   string
	method    string
	data      any
	called    bool
}

func (s *serverListerInvokerStub) Invoke(
	_ context.Context,
	actorType, actorID, method string,
	data any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	s.called = true
	s.actorType = actorType
	s.actorID = actorID
	s.method = method
	s.data = data
	return s.envelope, s.err
}

type serverListEnvelope struct {
	servers []ookla.Server
	err     error
}

func (e serverListEnvelope) Decode(into any) error {
	if e.err != nil {
		return e.err
	}
	*into.(*[]ookla.Server) = e.servers
	return nil
}

func testContext(t *testing.T) *echo.Context {
	t.Helper()
	return echo.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
}
