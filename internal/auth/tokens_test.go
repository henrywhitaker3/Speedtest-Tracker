package auth

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

func TestTokenRepositoryValidate(t *testing.T) {
	repo := testTokenRepository(t)
	ctx := context.Background()

	secret, _, err := repo.Create(ctx, "CI token")
	require.NoError(t, err)
	assert.NoError(t, repo.Validate(ctx, secret))

	err = repo.Validate(ctx, "unknown-token")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRegisterCreatesTokenHashIndex(t *testing.T) {
	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	require.NoError(t, err)
	require.NoError(t, Register(db))
	require.True(t, db.Migrator().HasIndex(&Token{}, "idx_tokens_hash"))

	first := Token{ID: uuid.MustOrdered(), Name: "first", Hash: "duplicate-hash"}
	require.NoError(t, db.Create(&first).Error)
	duplicate := Token{ID: uuid.MustOrdered(), Name: "duplicate", Hash: first.Hash}
	assert.Error(t, db.Create(&duplicate).Error)
}

func testTokenRepository(t *testing.T) *TokenRepository {
	t.Helper()

	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	require.NoError(t, err)
	require.NoError(t, Register(db))
	return NewRepo(db)
}
