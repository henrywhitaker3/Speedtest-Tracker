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

	result.Ping = pointer(12.5)
	result.Download = pointer(123_456.7)
	result.Upload = pointer(98_765.4)
	result.ServerID = pointer(42)
	result.ServerHost = pointer("speedtest.example.com")
	result.ServerName = pointer("Example Server")
	result.URL = pointer("https://www.speedtest.net/result/42")
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

	result.Ping = pointer(0.0)
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

func pointer[T any](value T) *T {
	return &value
}
