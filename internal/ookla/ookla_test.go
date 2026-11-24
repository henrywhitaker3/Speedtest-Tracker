package ookla

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/testutils"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
	"github.com/prometheus/client_golang/prometheus"
	metricspb "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

func TestDownloadURLForPlatform(t *testing.T) {
	const version = "1.2.0"

	tests := []struct {
		name   string
		goos   string
		goarch string
		want   string
	}{
		{
			name:   "macOS universal",
			goos:   "darwin",
			goarch: "arm64",
			want:   "https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-macosx-universal.tgz",
		},
		{
			name:   "Linux i386",
			goos:   "linux",
			goarch: "386",
			want:   "https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-i386.tgz",
		},
		{
			name:   "Linux x86_64",
			goos:   "linux",
			goarch: "amd64",
			want:   "https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-x86_64.tgz",
		},
		{
			name:   "Linux armhf",
			goos:   "linux",
			goarch: "arm",
			want:   "https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-armhf.tgz",
		},
		{
			name:   "Linux aarch64",
			goos:   "linux",
			goarch: "arm64",
			want:   "https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-aarch64.tgz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := downloadURLForPlatform(version, tt.goos, tt.goarch)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDownloadURLForPlatformRejectsUnsupportedPlatform(t *testing.T) {
	_, err := downloadURLForPlatform("1.2.0", "windows", "amd64")
	require.Error(t, err)
}

func TestExtractBinary(t *testing.T) {
	const contents = "speedtest executable"

	var archive bytes.Buffer
	gzipWriter := gzip.NewWriter(&archive)
	tarWriter := tar.NewWriter(gzipWriter)
	require.NoError(t, tarWriter.WriteHeader(
		&tar.Header{Name: "speedtest", Mode: 0o755, Size: int64(len(contents))},
	))
	_, err := tarWriter.Write([]byte(contents))
	require.NoError(t, err)
	require.NoError(t, tarWriter.Close())
	require.NoError(t, gzipWriter.Close())

	destination := filepath.Join(t.TempDir(), "missing", "speedtest")
	require.NoError(t, extractBinary(&archive, destination))

	got, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contents, string(got))

	info, err := os.Stat(destination)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
}

func TestOoklaOutputUnmarshalsFractionalPacketLoss(t *testing.T) {
	const packetLoss = 0.43859649122807015

	var output results.SpeedtestOutput
	require.NoError(t, json.Unmarshal([]byte(`{"packetLoss":0.43859649122807015}`), &output))
	assert.Equal(t, packetLoss, output.PacketLoss)
}

func TestPickServerSelectsConfiguredServer(t *testing.T) {
	serverIDs := []int{101, 202, 303}
	ookla := Ookla{server: serverIDs}

	picked := []int{}
	for range 100 {
		serverID, err := ookla.PickServer(t.Context())
		require.NoError(t, err)
		assert.Contains(t, serverIDs, serverID)
		if !slices.Contains(picked, serverID) {
			picked = append(picked, serverID)
		}
	}
	require.NotEqual(t, 1, len(picked))
}

func TestListServersCachesFreshServerList(t *testing.T) {
	servers := []Server{{ID: 1, Name: "nearest"}}
	caller := &callerStub{
		serverLists: []serverEnvelope{{
			Timestamp: time.Now().Add(-29 * time.Minute),
			Servers:   servers,
		}},
	}
	host := serverListerHost(t, caller)
	ookla := &Ookla{client: host.Service()}
	otherInstance := &Ookla{client: host.Service()}

	first, err := ookla.ListServers(t.Context())
	require.NoError(t, err)
	second, err := otherInstance.ListServers(t.Context())
	require.NoError(t, err)

	assert.Equal(t, servers, first)
	assert.Equal(t, servers, second)
	assert.Equal(t, 1, caller.listServersCalls)
}

func TestListServersRefreshesExpiredServerList(t *testing.T) {
	expiredServers := []Server{{ID: 1, Name: "expired"}}
	refreshedServers := []Server{{ID: 2, Name: "refreshed"}}
	caller := &callerStub{
		serverLists: []serverEnvelope{
			{Timestamp: time.Now().Add(-31 * time.Minute), Servers: expiredServers},
			{Timestamp: time.Now(), Servers: refreshedServers},
		},
	}
	ookla := &Ookla{client: serverListerHost(t, caller).Service()}

	first, err := ookla.ListServers(t.Context())
	require.NoError(t, err)
	second, err := ookla.ListServers(t.Context())
	require.NoError(t, err)

	assert.Equal(t, expiredServers, first)
	assert.Equal(t, refreshedServers, second)
	assert.Equal(t, 2, caller.listServersCalls)
}

func TestListServersDoesNotCacheFailedRefreshes(t *testing.T) {
	want := errors.New("list servers")
	servers := []Server{{ID: 1, Name: "nearest"}}
	caller := &callerStub{
		listServersErrors: []error{want},
		serverLists: []serverEnvelope{
			{},
			{Timestamp: time.Now(), Servers: servers},
		},
	}
	ookla := &Ookla{client: serverListerHost(t, caller).Service()}

	_, err := ookla.ListServers(t.Context())
	require.ErrorIs(t, err, want)
	got, err := ookla.ListServers(t.Context())
	require.NoError(t, err)

	assert.Equal(t, servers, got)
	assert.Equal(t, 2, caller.listServersCalls)
}

func TestRunUpdatesResult(t *testing.T) {
	output := speedtestOutput(10)
	ookla := &Ookla{caller: &callerStub{speedtestOutput: output}}
	result := &results.Result{}

	require.NoError(t, ookla.Run(t.Context(), output.Server.ID, result))

	assert.Equal(t, output.Ping.Latency, *result.Ping)
	assert.Equal(t, int64(output.Download.Bandwidth), *result.Download)
	assert.Equal(t, int64(output.Upload.Bandwidth), *result.Upload)
	assert.Equal(t, output.Server.ID, *result.ServerID)
	assert.Equal(t, "speedtest.example:8080", *result.ServerHost)
	assert.Equal(t, output.Server.Name, *result.ServerName)
	assert.Equal(t, output.Result.URL, *result.URL)
	assert.True(t, result.Complete)
	require.NotNil(t, result.Raw)
	assert.Equal(t, output, result.Raw.Data())
}

func TestRunUpdatesPrometheusMetrics(t *testing.T) {
	output := speedtestOutput(987654321)
	ookla := &Ookla{caller: &callerStub{speedtestOutput: output}}
	labels := []string{"987654321"}
	runsBefore := metricValue(t, speedtestsRun.WithLabelValues(labels...)).GetCounter().GetValue()
	durationCountBefore := metricValue(
		t,
		metricWithLabels(t, speedtestDurationSeconds.MetricVec, labels...),
	).GetHistogram().
		GetSampleCount()

	require.NoError(t, ookla.Run(t.Context(), output.Server.ID, &results.Result{}))

	assert.Equal(
		t,
		runsBefore+1,
		metricValue(t, speedtestsRun.WithLabelValues(labels...)).GetCounter().GetValue(),
	)
	assert.Equal(
		t,
		durationCountBefore+1,
		metricValue(
			t,
			metricWithLabels(t, speedtestDurationSeconds.MetricVec, labels...),
		).GetHistogram().
			GetSampleCount(),
	)
	assert.Equal(
		t,
		output.Ping.Latency,
		metricValue(t, speedtestPingMilliseconds.WithLabelValues(labels...)).GetGauge().GetValue(),
	)
	assert.Equal(
		t,
		float64(output.Download.Bandwidth),
		metricValue(
			t,
			speedtestDownloadBytesPerSecond.WithLabelValues(labels...),
		).GetGauge().
			GetValue(),
	)
	assert.Equal(
		t,
		float64(output.Upload.Bandwidth),
		metricValue(
			t,
			speedtestUploadBytesPerSecond.WithLabelValues(labels...),
		).GetGauge().
			GetValue(),
	)
}

func speedtestOutput(serverID int) results.SpeedtestOutput {
	return results.SpeedtestOutput{
		Ping:     results.Ping{Latency: 10.5},
		Download: results.Measurement{Bandwidth: 1_000_000},
		Upload:   results.Measurement{Bandwidth: 500_000},
		Server: &results.Server{
			ID:   serverID,
			Host: "speedtest.example",
			Port: 8080,
			Name: "Speedtest",
		},
		Result: &results.OoklaResult{URL: "https://www.speedtest.net/result/1"},
	}
}

func metricValue(t *testing.T, metric prometheus.Metric) *metricspb.Metric {
	t.Helper()
	value := &metricspb.Metric{}
	require.NoError(t, metric.Write(value))
	return value
}

func metricWithLabels(
	t *testing.T,
	vector *prometheus.MetricVec,
	labels ...string,
) prometheus.Metric {
	t.Helper()
	metric, err := vector.GetMetricWithLabelValues(labels...)
	require.NoError(t, err)
	return metric
}

type callerStub struct {
	serverLists       []serverEnvelope
	listServersErrors []error
	listServersCalls  int
	speedtestOutput   results.SpeedtestOutput
}

func (c *callerStub) listServers(context.Context) (serverEnvelope, error) {
	call := c.listServersCalls
	c.listServersCalls++

	if call < len(c.listServersErrors) && c.listServersErrors[call] != nil {
		return serverEnvelope{}, c.listServersErrors[call]
	}
	if call < len(c.serverLists) {
		return c.serverLists[call], nil
	}
	return serverEnvelope{}, nil
}

func serverListerHost(t *testing.T, lister lister) *local.Host {
	t.Helper()
	return testutils.Host(t, func(t testing.TB, host *local.Host) {
		require.NoError(t, host.RegisterSingletonActor(
			Actor,
			func(actorID string, service *actor.Service) actor.Actor {
				return &serverListerActor{
					client: actor.NewActorClient[serverListerState](Actor, actorID, service),
					lister: lister,
				}
			},
		))
	})
}

func (c *callerStub) runTest(context.Context, int) (results.SpeedtestOutput, error) {
	return c.speedtestOutput, nil
}
