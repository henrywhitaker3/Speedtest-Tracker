package ookla

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("downloadURLForPlatform() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDownloadURLForPlatformRejectsUnsupportedPlatform(t *testing.T) {
	if _, err := downloadURLForPlatform("1.2.0", "windows", "amd64"); err == nil {
		t.Fatal("downloadURLForPlatform() succeeded for an unsupported platform")
	}
}

func TestExtractBinary(t *testing.T) {
	const contents = "speedtest executable"

	var archive bytes.Buffer
	gzipWriter := gzip.NewWriter(&archive)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := tarWriter.WriteHeader(
		&tar.Header{Name: "speedtest", Mode: 0o755, Size: int64(len(contents))},
	); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write([]byte(contents)); err != nil {
		t.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}

	destination := filepath.Join(t.TempDir(), "missing", "speedtest")
	if err := extractBinary(&archive, destination); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != contents {
		t.Errorf("extracted contents = %q, want %q", got, contents)
	}

	info, err := os.Stat(destination)
	if err != nil {
		t.Fatal(err)
	}
	if permissions := info.Mode().Perm(); permissions != 0o755 {
		t.Errorf("permissions = %o, want 755", permissions)
	}
}

func TestOoklaOutputUnmarshalsFractionalPacketLoss(t *testing.T) {
	const packetLoss = 0.43859649122807015

	var output results.SpeedtestOutput
	if err := json.Unmarshal([]byte(`{"packetLoss":0.43859649122807015}`), &output); err != nil {
		t.Fatal(err)
	}
	if output.PacketLoss != packetLoss {
		t.Errorf("packet loss = %v, want %v", output.PacketLoss, packetLoss)
	}
}
