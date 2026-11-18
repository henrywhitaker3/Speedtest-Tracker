package ookla

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"go.opentelemetry.io/otel/attribute"
)

func NewCaller() *ooklaCaller {
	return &ooklaCaller{}
}

type ooklaCaller struct{}

func (o *ooklaCaller) listServers(ctx context.Context) (serverEnvelope, error) {
	execArgs := []string{"--servers", "--format=json"}
	execCtx, span := tracing.NewSpan(ctx, "exec")
	span.SetAttributes(
		attribute.String("exec.binary", binaryFile),
		attribute.StringSlice("exec.args", execArgs),
	)

	cmd := exec.CommandContext(execCtx, binaryFile, execArgs...)
	raw, err := cmd.Output()
	span.End()
	if err != nil {
		err := fmt.Errorf("list nearest servers: %w", err)
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			err = fmt.Errorf("%w: %s", err, string(exitErr.Stderr))
		}
		return serverEnvelope{}, err
	}

	var servers serverEnvelope
	if err := json.Unmarshal(raw, &servers); err != nil {
		return serverEnvelope{}, fmt.Errorf("unmarshal server data: %w", err)
	}

	return servers, nil
}

func (o *ooklaCaller) runTest(ctx context.Context, server int) (results.SpeedtestOutput, error) {
	execArgs := append([]string{}, agrs...)
	if server != 0 {
		execArgs = append(execArgs, fmt.Sprintf("--server-id=%d", server))
	}

	execCtx, execSpan := tracing.NewSpan(ctx, "exec")
	execSpan.SetAttributes(
		attribute.String("exec.binary", binaryFile),
		attribute.StringSlice("exec.args", execArgs),
	)

	cmd := exec.CommandContext(execCtx, binaryFile, execArgs...)
	raw, err := cmd.Output()
	execSpan.End()
	if err != nil {
		speedtestErrors.Inc()
		err := fmt.Errorf("run speedtest: %w", err)
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			err = fmt.Errorf("%w: %s", err, string(exitErr.Stderr))
		}
		return results.SpeedtestOutput{}, err
	}

	slog.DebugContext(ctx, "got raw speedtest output", "output", string(raw))
	var output results.SpeedtestOutput
	if err := json.Unmarshal(raw, &output); err != nil {
		speedtestErrors.Inc()
		return results.SpeedtestOutput{}, fmt.Errorf("unmarshal output: %w", err)
	}
	slog.DebugContext(ctx, "got parsed speedtest output", "output", output)

	return output, nil
}
