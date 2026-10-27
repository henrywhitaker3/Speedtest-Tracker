package notifications

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/units"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

type telegramMessageSender interface {
	SendMessage(context.Context, *bot.SendMessageParams) (*models.Message, error)
}

func SendTelegram(
	ctx context.Context,
	sender telegramMessageSender,
	n config.Notification,
	uiURL string,
	phase Phase,
	result results.Result,
) error {
	var msg string
	switch phase {
	case PreRun:
		msg = buildTestStartingMessage(result, uiURL)
	case PostRun:
		if result.Error != nil {
			msg = buildErrorMessage(result, uiURL)
		} else {
			msg = buildSuccessMessage(result, uiURL)
		}
	default:
		return fmt.Errorf("unknown phase %d", phase)
	}

	params := &bot.SendMessageParams{
		ChatID:    string(n.ChatID),
		Text:      msg,
		ParseMode: models.ParseModeMarkdownV1,
	}

	if _, err := sender.SendMessage(ctx, params); err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}

func buildTestStartingMessage(r results.Result, uiURL string) string {
	return fmt.Sprintf("*Starting speedtest*\n%s", buildMarkdownURL(uiURL, r))
}

func buildSuccessMessage(r results.Result, uiURL string) string {
	return fmt.Sprintf(
		"*New speedtest*\nPing: *%.2f ms*\nDownload: *%s*\nUpload: *%s*\n%s",
		*r.Ping,
		units.BitsHuman(units.BytesToBits(*r.Download)),
		units.BitsHuman(units.BytesToBits(*r.Upload)),
		buildMarkdownURL(uiURL, r),
	)
}

func buildErrorMessage(r results.Result, uiURL string) string {
	return fmt.Sprintf(
		"*Failed speedtest*\nError: %s\n%s",
		*r.Error,
		buildMarkdownURL(uiURL, r),
	)
}

func buildMarkdownURL(uiURL string, r results.Result) string {
	return fmt.Sprintf("[See result](%s)", buildResultsURL(uiURL, r))
}

func buildResultsURL(uiURL string, r results.Result) string {
	return fmt.Sprintf("%s/results/%s", uiURL, r.ID.String())
}
