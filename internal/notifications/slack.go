package notifications

import (
	"context"
	"fmt"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/units"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/slack-go/slack"
)

func sendSlack(ctx context.Context, url, uiURL string, phase Phase, result results.Result) error {
	var msg *slack.WebhookMessage
	switch phase {
	case PreRun:
		msg = buildSlackPreRun(result, uiURL)
	case PostRun:
		msg = buildSlackPostRun(result, uiURL)
	default:
		return fmt.Errorf("unknown phase %d", phase)
	}

	return slack.PostWebhookContext(ctx, url, msg)
}

func buildSlackPreRun(r results.Result, uiURL string) *slack.WebhookMessage {
	return &slack.WebhookMessage{
		Text: StartingTitle,
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(
					slack.NewTextBlockObject("plain_text", StartingTitle, true, false),
				),
				slack.NewSectionBlock(
					buildSlackResultButton(r, uiURL),
					nil,
					nil,
				),
			},
		},
	}
}

func buildSlackPostRun(r results.Result, uiURL string) *slack.WebhookMessage {
	if r.Error != nil {
		return buildSlackErrorMessage(r, uiURL)
	}
	return buildSlackSuccessMessage(r, uiURL)
}

func buildSlackSuccessMessage(r results.Result, uiURL string) *slack.WebhookMessage {
	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Ping*\n%.2f ms", *r.Ping), false, false),
		slack.NewTextBlockObject(
			"mrkdwn",
			fmt.Sprintf("*Download*\n%s", units.BitsHuman(units.BytesToBits(*r.Download))),
			false,
			false,
		),
		slack.NewTextBlockObject(
			"mrkdwn",
			fmt.Sprintf("*Upload*\n%s", units.BitsHuman(units.BytesToBits(*r.Upload))),
			false,
			false,
		),
	}

	return &slack.WebhookMessage{
		Text: SuccessTitle,
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(
					slack.NewTextBlockObject("plain_text", SuccessTitle, true, false),
				),
				slack.NewSectionBlock(nil, fields, nil),
				slack.NewSectionBlock(
					buildSlackResultButton(r, uiURL),
					nil,
					nil,
				),
			},
		},
	}
}

func buildSlackErrorMessage(r results.Result, uiURL string) *slack.WebhookMessage {
	return &slack.WebhookMessage{
		Text: ErrorTitle,
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(
					slack.NewTextBlockObject("plain_text", ErrorTitle, true, false),
				),
				slack.NewSectionBlock(
					slack.NewTextBlockObject(
						"mrkdwn",
						fmt.Sprintf("*Error*\n%s", *r.Error),
						false,
						false,
					),
					nil,
					nil,
				),
				slack.NewSectionBlock(
					buildSlackResultButton(r, uiURL),
					nil,
					nil,
				),
			},
		},
	}
}

func buildSlackResultButton(r results.Result, uiURL string) *slack.TextBlockObject {
	return slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("<%s|View Result>", buildResultsURL(uiURL, r)),
		false,
		false,
	)
}
