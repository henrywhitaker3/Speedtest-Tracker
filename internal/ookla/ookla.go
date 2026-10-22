// Package ookla
package ookla

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
)

func Run(ctx context.Context, server int, result *results.Result) error {
	ctx, span := tracing.NewSpan(ctx, "Run speedtest")
	defer span.End()

	execArgs := append([]string{}, agrs...)
	if server != 0 {
		execArgs = append(execArgs, fmt.Sprintf("--server-id=%d", server))
	}

	cmd := exec.CommandContext(ctx, binaryFile, execArgs...)
	raw, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("run speedtest: %w", err)
	}

	slog.DebugContext(ctx, "got raw speedtest output", "output", string(raw))
	var output ooklaOutput
	if err := json.Unmarshal(raw, &output); err != nil {
		return fmt.Errorf("unmarshal output: %w", err)
	}
	slog.DebugContext(ctx, "got parsed speedtest output", "output", output)

	result.Ping = new(output.Ping.Latency)
	result.Download = new(convertResult(output.Download.Bandwidth))
	result.Upload = new(convertResult(output.Upload.Bandwidth))
	result.ServerID = new(output.Server.ID)
	result.ServerName = new(output.Server.Name)
	result.ServerHost = new(fmt.Sprintf("%s:%d", output.Server.Host, output.Server.Port))
	result.URL = new(output.Result.URL)
	result.Complete = true

	return nil
}

func convertResult(raw int) float64 {
	return (float64(raw) * 8) / 1000000
}

const (
	binaryFile          = "/tmp/speedtest-cli"
	downloadURLTemplate = "https://install.speedtest.net/app/cli/ookla-speedtest-%s-%s.tgz"
)

var (
	agrs = []string{
		"--format=json",
		"--accept-license",
		"--accept-gdpr",
	}
	linuxArchitectures = map[string]string{
		"386":   "i386",
		"amd64": "x86_64",
		"arm":   "armhf",
		"arm64": "aarch64",
	}
)

func downloadURL(version string) (string, error) {
	return downloadURLForPlatform(version, runtime.GOOS, runtime.GOARCH)
}

func downloadURLForPlatform(version, goos, goarch string) (string, error) {
	switch goos {
	case "darwin":
		return fmt.Sprintf(downloadURLTemplate, version, "macosx-universal"), nil
	case "linux":
		architecture, ok := linuxArchitectures[goarch]
		if !ok {
			return "", fmt.Errorf("unsupported Linux architecture %q", goarch)
		}
		return fmt.Sprintf(downloadURLTemplate, version, "linux-"+architecture), nil
	default:
		return "", fmt.Errorf("unsupported operating system %q", goos)
	}
}

// EnsureBinary downloads the ookla speedtest cli
func EnsureBinary(ctx context.Context, version string) error {
	url, err := downloadURL(version)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "downloading ookla binary", "url", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}

	if err := extractBinary(resp.Body, binaryFile); err != nil {
		return fmt.Errorf("extract binary: %w", err)
	}

	return nil
}

func extractBinary(archive io.Reader, destination string) error {
	gzipReader, err := gzip.NewReader(archive)
	if err != nil {
		return fmt.Errorf("open gzip stream: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return fmt.Errorf("speedtest executable not found")
		}
		if err != nil {
			return fmt.Errorf("read archive: %w", err)
		}
		if !header.FileInfo().Mode().IsRegular() || path.Base(header.Name) != "speedtest" {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create binary directory: %w", err)
		}
		binary, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return fmt.Errorf("open binary file: %w", err)
		}
		if _, err := io.Copy(binary, tarReader); err != nil {
			binary.Close()
			return fmt.Errorf("write binary: %w", err)
		}
		if err := binary.Chmod(0o755); err != nil {
			binary.Close()
			return fmt.Errorf("set binary permissions: %w", err)
		}
		if err := binary.Close(); err != nil {
			return fmt.Errorf("close binary file: %w", err)
		}

		return nil
	}
}

type ooklaOutput struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
	} `json:"ping"`
	Download struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"upload"`
	PacketLoss int    `json:"packetLoss"`
	Isp        string `json:"isp"`
	Interface  struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}
