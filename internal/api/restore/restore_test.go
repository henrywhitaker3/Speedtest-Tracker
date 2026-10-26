package restore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/labstack/echo/v5"
)

func TestRestoreHandlerSkipsResultsWithZeroCreatedAt(t *testing.T) {
	createdAt := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	repo := &restoreRepoStub{restored: 1}
	handler := (&restoreHandler{repo: repo}).Handler()

	got, err := handler(testContext(t), []LegacyResult{
		{CreatedAt: LegcayTime{}},
		{CreatedAt: LegcayTime(createdAt)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Restored != 1 || got.Skipped != 1 {
		t.Errorf("Handler() response = %+v, want restored 1 and skipped 1", got)
	}
	if len(repo.results) != 1 {
		t.Fatalf("Restore() received %d results, want 1", len(repo.results))
	}
	if !repo.results[0].CreatedAt.Equal(createdAt) {
		t.Errorf("Restore() received created_at %v, want %v", repo.results[0].CreatedAt, createdAt)
	}
}

type restoreRepoStub struct {
	restored int
	skipped  int
	err      error
	results  []results.Result
}

func (r *restoreRepoStub) Restore(_ context.Context, got []results.Result) (int, int, error) {
	r.results = got
	return r.restored, r.skipped, r.err
}

func testContext(t *testing.T) *echo.Context {
	t.Helper()
	return echo.NewContext(httptest.NewRequest(http.MethodPost, "/v1/restore/json", nil), httptest.NewRecorder())
}
