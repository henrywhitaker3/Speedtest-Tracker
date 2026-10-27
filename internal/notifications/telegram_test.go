package notifications

import (
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestTelegramMessageBuildersUseMarkdownResultLinks(t *testing.T) {
	result := results.Result{ID: uuid.MustOrdered()}
	resultURL := "https://speedtest.example/results/" + result.ID.String()

	require.Equal(t, resultURL, buildResultsURL("https://speedtest.example", result))
	require.Equal(
		t,
		"[See result]("+resultURL+")",
		buildMarkdownURL("https://speedtest.example", result),
	)
}
