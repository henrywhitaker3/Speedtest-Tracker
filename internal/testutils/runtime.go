// Package testutils
package testutils

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/components/sqlite"
	"github.com/italypaleale/francis/host/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HostOption func(testing.TB, *local.Host)

func WithWorkflow(w *workflow.Workflow) HostOption {
	return func(t testing.TB, h *local.Host) {
		require.Nil(t, h.RegisterBuiltInActor(w))
	}
}

func WithActor(kind string, f actor.Factory) HostOption {
	return func(t testing.TB, h *local.Host) {
		require.Nil(t, h.RegisterActor(kind, f))
	}
}

func Host(t testing.TB, opts ...HostOption) *local.Host {
	t.Helper()

	psk := make([]byte, 64)
	_, err := rand.Read(psk)
	require.Nil(t, err)
	address := freeUDPAddress(t)
	host, err := local.NewHost(
		local.WithAddress(address),
		local.WithSQLiteProvider(sqlite.SQLiteProviderOptions{
			ConnectionString: filepath.Join(t.TempDir(), "workflow.db"),
		}),
		local.WithRuntimePSKs(psk),
		local.WithAlarmsPollInterval(100*time.Millisecond),
		local.WithLogger(slog.New(slog.DiscardHandler)),
	)
	require.NoError(t, err)

	for _, o := range opts {
		o(t, host)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- host.Run(ctx) }()
	select {
	case <-host.Ready():
	case <-time.After(10 * time.Second):
		cancel()
		require.FailNow(t, "assertion failed", "workflow host did not become ready")
	}
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-errCh:
			assert.Truef(
				t,
				err == nil || errors.Is(err, context.Canceled),
				"workflow host stopped with error: %v",
				err,
			)
		case <-time.After(10 * time.Second):
			assert.Fail(t, "workflow host did not stop")
		}
	})
	return host
}

func freeUDPAddress(t testing.TB) string {
	t.Helper()

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	require.NoError(t, err)
	address := listener.LocalAddr().String()
	require.NoError(t, listener.Close())
	return address
}
