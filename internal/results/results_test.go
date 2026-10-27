package results

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

func TestRepositoryCreateAndGet(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	created, err := repo.Create(ctx)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, created.ID)
	assert.False(t, created.Complete)

	got, err := repo.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.False(t, got.Complete)
}

func TestRepositoryUpdateAndFresh(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	result, err := repo.Create(ctx)
	require.NoError(t, err)

	result.Ping = new(12.5)
	result.Download = new(int64(123_456))
	result.Upload = new(int64(98_765))
	result.ServerID = new(42)
	result.ServerHost = new("speedtest.example.com")
	result.ServerName = new("Example Server")
	result.URL = new("https://www.speedtest.net/result/42")
	result.Complete = true
	require.NoError(t, repo.Update(ctx, result))

	stored, err := repo.Get(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, stored.Ping)
	assert.Equal(t, 12.5, *stored.Ping)
	require.NotNil(t, stored.Download)
	assert.Equal(t, int64(123_456), *stored.Download)
	require.NotNil(t, stored.ServerName)
	assert.Equal(t, "Example Server", *stored.ServerName)
	assert.True(t, stored.Complete)

	result.Ping = new(0.0)
	result.Complete = false
	require.NoError(t, repo.Fresh(ctx, result))
	require.NotNil(t, result.Ping)
	assert.Equal(t, 12.5, *result.Ping)
	assert.True(t, result.Complete)
}

func TestRepositoryGetMissing(t *testing.T) {
	repo := testRepository(t)

	_, err := repo.Get(context.Background(), uuid.MustOrdered())
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepositoryListPaginatesNewestFirst(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	created := make([]*Result, 3)
	for i := range created {
		result, err := repo.Create(ctx)
		require.NoError(t, err)
		created[i] = result
	}

	firstPage, err := repo.List(ctx, ListOpts{Page: 1, Limit: 2})
	require.NoError(t, err)
	assert.Equal(t, 1, firstPage.Page)
	assert.Equal(t, 2, firstPage.Limit)
	assert.Equal(t, int64(3), firstPage.Total)
	assert.Equal(t, 2, firstPage.TotalPages)
	require.Len(t, firstPage.Data, 2)
	assert.Equal(t, created[2].ID, firstPage.Data[0].ID)
	assert.Equal(t, created[1].ID, firstPage.Data[1].ID)

	secondPage, err := repo.List(ctx, ListOpts{Page: 2, Limit: 2})
	require.NoError(t, err)
	require.Len(t, secondPage.Data, 1)
	assert.Equal(t, created[0].ID, secondPage.Data[0].ID)
}

func TestRepositoryListUsesDefaultPaginationForEmptyResults(t *testing.T) {
	repo := testRepository(t)

	page, err := repo.List(context.Background(), ListOpts{})
	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 25, page.Limit)
	assert.Zero(t, page.Total)
	assert.Zero(t, page.TotalPages)
	assert.Empty(t, page.Data)
}

func TestRepositoryTimeseriesQueriesResults(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	complete, err := repo.Create(ctx)
	require.NoError(t, err)
	complete.Ping = new(12.5)
	complete.Download = new(int64(123_456))
	complete.Upload = new(int64(98_765))
	complete.Complete = true
	require.NoError(t, repo.Update(ctx, complete))

	incomplete, err := repo.Create(ctx)
	require.NoError(t, err)

	stats, err := repo.Timeseries(ctx, time.Hour)
	require.NoError(t, err)
	require.Len(t, stats, 1)
	assert.Equal(t, complete.ID, stats[0].ID)
	require.NotNil(t, stats[0].Ping)
	assert.Equal(t, *complete.Ping, *stats[0].Ping)
	assert.False(t, incomplete.Complete)
}

func TestRepositoryForgetDeletesExpiredResults(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	expired, err := repo.Create(ctx)
	require.NoError(t, err)
	current, err := repo.Create(ctx)
	require.NoError(t, err)

	require.NoError(t, repo.db.Model(&Result{}).
		Where("id = ?", expired.ID).
		Update("created_at", time.Now().Add(-2*time.Hour)).Error)
	require.NoError(t, repo.db.Model(&Result{}).
		Where("id = ?", current.ID).
		Update("created_at", time.Now().Add(-30*time.Minute)).Error)

	deleted, err := repo.Forget(ctx, time.Hour)
	require.NoError(t, err)
	assert.Equal(t, 1, deleted)

	var deletedResult Result
	err = repo.db.Unscoped().Where("id = ?", expired.ID).First(&deletedResult).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	_, err = repo.Get(ctx, current.ID)
	assert.NoError(t, err)
}

func TestRepositoryForgetWithZeroPeriodDeletesAllResults(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	for range 2 {
		_, err := repo.Create(ctx)
		require.NoError(t, err)
	}

	deleted, err := repo.Forget(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, deleted)

	var remaining int64
	require.NoError(t, repo.db.Unscoped().Model(&Result{}).Count(&remaining).Error)
	assert.Zero(t, remaining)
}

func TestRepositoryRestoreBatchesAndSkipsDuplicateCreatedAt(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	first := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	existing := restoreResult(first)
	require.NoError(t, repo.db.Create(&existing).Error)

	toRestore := []Result{restoreResult(first)}
	for i := 1; i <= restoreBatchSize+1; i++ {
		toRestore = append(toRestore, restoreResult(first.Add(time.Duration(i)*time.Second)))
	}
	toRestore = append(toRestore, restoreResult(first.Add(time.Second)))

	restored, skipped, err := repo.Restore(ctx, toRestore)
	require.NoError(t, err)
	assert.Equal(t, restoreBatchSize+1, restored)
	assert.Equal(t, 2, skipped)

	var count int64
	require.NoError(t, repo.db.Model(&Result{}).Count(&count).Error)
	assert.EqualValues(t, restoreBatchSize+2, count)
}

func TestRepositoryRestoreSkipsDuplicateCreatedAtWithinBatch(t *testing.T) {
	repo := testRepository(t)
	ctx := context.Background()

	createdAt := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	restored, skipped, err := repo.Restore(ctx, []Result{
		restoreResult(createdAt),
		restoreResult(createdAt),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, restored)
	assert.Equal(t, 1, skipped)

	var count int64
	require.NoError(t, repo.db.Model(&Result{}).Count(&count).Error)
	assert.EqualValues(t, 1, count)
}

func restoreResult(createdAt time.Time) Result {
	return Result{
		ID:        uuid.Must(uuid.OrderedAt(createdAt)),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Complete:  true,
	}
}

func TestRegisterCreatesTimeseriesIndex(t *testing.T) {
	db, err := database.Sqlite(filepath.Join(t.TempDir(), "results.db"))
	require.NoError(t, err)
	require.NoError(t, Register(db))
	require.True(t, db.Migrator().HasIndex(&Result{}, "idx_results_timeseries"))

	var definition string
	require.NoError(t, db.Raw(
		"SELECT sql FROM sqlite_master WHERE type = 'index' AND name = ?",
		"idx_results_timeseries",
	).Scan(&definition).Error)
	assert.Contains(t, definition, "created_at")
	assert.Contains(t, definition, "complete = true AND deleted_at IS NULL")
}

func testRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := database.Sqlite(filepath.Join(t.TempDir(), "results.db"))
	require.NoError(t, err)
	require.NoError(t, Register(db))

	return NewRepo(db)
}
