package ookla

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

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
