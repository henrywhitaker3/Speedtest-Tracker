package notifications

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/units"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

func sendGotify(
	ctx context.Context,
	endpoint, token, uiURL string,
	phase Phase,
	result results.Result,
) error {
	gotifyURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse gotify url: %w", err)
	}

	payload, err := buildGotifyMessage(uiURL, phase, result)
	if err != nil {
		return err
	}

	params := message.NewCreateMessageParamsWithContext(ctx).WithBody(payload)
	if _, err := gotify.NewClient(gotifyURL, http.DefaultClient).Message.CreateMessage(
		params,
		auth.TokenAuth(token),
	); err != nil {
		return fmt.Errorf("send gotify message: %w", err)
	}

	return nil
}

func buildGotifyMessage(
	uiURL string,
	phase Phase,
	result results.Result,
) (*models.MessageExternal, error) {
	message := &models.MessageExternal{
		Extras: map[string]any{
			"client::display": map[string]string{"contentType": "text/markdown"},
			"client::notification": map[string]any{
				"click": map[string]string{"url": buildResultsURL(uiURL, result)},
			},
		},
	}

	switch phase {
	case PreRun:
		message.Title = StartingTitle
		message.Message = buildMarkdownURL(uiURL, result)
		message.Priority = 5
	case PostRun:
		if result.Error != nil {
			message.Title = ErrorTitle
			message.Message = fmt.Sprintf(
				"Error: %s\n%s",
				*result.Error,
				buildMarkdownURL(uiURL, result),
			)
			message.Priority = 10
			return message, nil
		}

		message.Title = SuccessTitle
		message.Message = fmt.Sprintf(
			"Ping: **%.2f ms**\nDownload: **%s**\nUpload: **%s**\n%s",
			*result.Ping,
			units.BitsHuman(units.BytesToBits(*result.Download)),
			units.BitsHuman(units.BytesToBits(*result.Upload)),
			buildMarkdownURL(uiURL, result),
		)
		message.Priority = 5
	default:
		return nil, fmt.Errorf("unknown phase %d", phase)
	}

	return message, nil
}
