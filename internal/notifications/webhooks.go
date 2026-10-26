package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

func notifyHealthcheck(ctx context.Context, url string, phase Phase, result results.Result) error {
	switch phase {
	case PreRun:
		url += "/start"
	case PostRun:
		if result.Error != nil {
			url += "/fail"
		}
	}
	return sendWebhook(ctx, url, result)
}

func sendWebhook(ctx context.Context, url string, result results.Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bodyReader,
	)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}

	return nil
}
