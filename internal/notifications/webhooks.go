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
	return sendWebhook(ctx, url, nil, result)
}

func sendWebhook(
	ctx context.Context,
	url string,
	headers map[string]string,
	result results.Result,
) error {
	return sendJSON(ctx, url, result, headers)
}

func sendJSON(ctx context.Context, url string, payload any, headers map[string]string) error {
	body, err := json.Marshal(payload)
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
	for key, val := range headers {
		req.Header.Set(key, val)
	}

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
