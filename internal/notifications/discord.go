package notifications

import (
	"context"
	"fmt"

	discord "github.com/bwmarrin/discordgo"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/units"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

const (
	discordInfoColor    = 0x3498DB
	discordSuccessColor = 0x2ECC71
	discordFailureColor = 0xE74C3C
)

func sendDiscord(ctx context.Context, url, uiURL string, phase Phase, result results.Result) error {
	var params *discord.WebhookParams
	switch phase {
	case PreRun:
		params = buildDiscordPreRun(result, uiURL)
	case PostRun:
		params = buildDiscordPostRun(result, uiURL)
	default:
		return fmt.Errorf("unknown phase %d", phase)
	}

	if err := sendJSON(ctx, url, params, nil); err != nil {
		return fmt.Errorf("send discord webhook: %w", err)
	}
	return nil
}

func buildDiscordPreRun(r results.Result, uiURL string) *discord.WebhookParams {
	return &discord.WebhookParams{
		Embeds: []*discord.MessageEmbed{{
			Title: StartingTitle,
			URL:   buildResultsURL(uiURL, r),
			Color: discordInfoColor,
		}},
	}
}

func buildDiscordPostRun(r results.Result, uiURL string) *discord.WebhookParams {
	if r.Error != nil {
		return &discord.WebhookParams{
			Embeds: []*discord.MessageEmbed{{
				Title: ErrorTitle,
				URL:   buildResultsURL(uiURL, r),
				Color: discordFailureColor,
				Fields: []*discord.MessageEmbedField{{
					Name:  "Error",
					Value: *r.Error,
				}},
			}},
		}
	}

	return &discord.WebhookParams{
		Embeds: []*discord.MessageEmbed{{
			Title: SuccessTitle,
			URL:   buildResultsURL(uiURL, r),
			Color: discordSuccessColor,
			Fields: []*discord.MessageEmbedField{
				{
					Name:   "Ping",
					Value:  fmt.Sprintf("%.2f ms", *r.Ping),
					Inline: true,
				},
				{
					Name:   "Download",
					Value:  units.BitsHuman(units.BytesToBits(*r.Download)),
					Inline: true,
				},
				{
					Name:   "Upload",
					Value:  units.BitsHuman(units.BytesToBits(*r.Upload)),
					Inline: true,
				},
			},
		}},
	}
}
