package results

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

func TestRepositoryCreateAndGet(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	created, err := repo.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == (uuid.UUID{}) {
		t.Fatal("created result has an empty ID")
	}
	if created.Complete {
		t.Fatal("new result is marked complete")
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("got ID %q, want %q", got.ID, created.ID)
	}
	if got.Complete {
		t.Fatal("stored result is marked complete")
	}
}

func TestRepositoryUpdateAndFresh(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	result, err := repo.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result.Ping = new(12.5)
	result.Download = new(123_456.7)
	result.Upload = new(98_765.4)
	result.ServerID = new(42)
	result.ServerHost = new("speedtest.example.com")
	result.ServerName = new("Example Server")
	result.URL = new("https://www.speedtest.net/result/42")
	result.Complete = true
	if err := repo.Update(ctx, result); err != nil {
		t.Fatal(err)
	}

	stored, err := repo.Get(ctx, result.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Ping == nil || *stored.Ping != 12.5 {
		t.Errorf("stored ping = %v, want 12.5", stored.Ping)
	}
	if stored.Download == nil || *stored.Download != 123_456.7 {
		t.Errorf("stored download = %v, want 123456.7", stored.Download)
	}
	if stored.ServerName == nil || *stored.ServerName != "Example Server" {
		t.Errorf("stored server name = %v, want Example Server", stored.ServerName)
	}
	if !stored.Complete {
		t.Fatal("stored result is not complete")
	}

	result.Ping = new(0.0)
	result.Complete = false
	if err := repo.Fresh(ctx, result); err != nil {
		t.Fatal(err)
	}
	if result.Ping == nil || *result.Ping != 12.5 {
		t.Errorf("fresh ping = %v, want 12.5", result.Ping)
	}
	if !result.Complete {
		t.Fatal("fresh result is not complete")
	}
}

func TestRepositoryGetMissing(t *testing.T) {
	repo := testRepository(t)

	_, err := repo.Get(context.Background(), uuid.MustOrdered())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Get() error = %v, want record-not-found error", err)
	}
}

func TestRepositoryListPaginatesNewestFirst(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	created := make([]*Result, 3)
	for i := range created {
		result, err := repo.Create(ctx)
		if err != nil {
			t.Fatal(err)
		}
		created[i] = result
	}

	firstPage, err := repo.List(ctx, ListOpts{Page: 1, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if firstPage.Page != 1 || firstPage.Limit != 2 {
		t.Errorf("first page metadata = page %d, limit %d; want page 1, limit 2", firstPage.Page, firstPage.Limit)
	}
	if firstPage.Total != 3 || firstPage.TotalPages != 2 {
		t.Errorf("first page totals = %d rows, %d pages; want 3 rows, 2 pages", firstPage.Total, firstPage.TotalPages)
	}
	if len(firstPage.Data) != 2 {
		t.Fatalf("first page has %d results, want 2", len(firstPage.Data))
	}
	if firstPage.Data[0].ID != created[2].ID || firstPage.Data[1].ID != created[1].ID {
		t.Errorf("first page IDs = [%q %q], want [%q %q]", firstPage.Data[0].ID, firstPage.Data[1].ID, created[2].ID, created[1].ID)
	}

	secondPage, err := repo.List(ctx, ListOpts{Page: 2, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(secondPage.Data) != 1 {
		t.Fatalf("second page has %d results, want 1", len(secondPage.Data))
	}
	if secondPage.Data[0].ID != created[0].ID {
		t.Errorf("second page ID = %q, want %q", secondPage.Data[0].ID, created[0].ID)
	}
}

func TestRepositoryListUsesDefaultPaginationForEmptyResults(t *testing.T) {
	repo := testRepository(t)

	page, err := repo.List(context.Background(), ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 1 || page.Limit != 25 {
		t.Errorf("default metadata = page %d, limit %d; want page 1, limit 25", page.Page, page.Limit)
	}
	if page.Total != 0 || page.TotalPages != 0 {
		t.Errorf("empty page totals = %d rows, %d pages; want 0 rows, 0 pages", page.Total, page.TotalPages)
	}
	if len(page.Data) != 0 {
		t.Errorf("empty page has %d results, want 0", len(page.Data))
	}
}

func testRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := database.Sqlite(filepath.Join(t.TempDir(), "results.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := Register(db); err != nil {
		t.Fatal(err)
	}

	return NewRepo(db)
}
