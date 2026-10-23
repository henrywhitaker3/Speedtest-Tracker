package auth

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

func TestTokenRepositoryValidate(t *testing.T) {
	repo := testTokenRepository(t)
	ctx := context.Background()

	secret, _, err := repo.Create(ctx, "CI token")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := repo.Validate(ctx, secret); err != nil {
		t.Errorf("Validate() error = %v, want nil for a stored token", err)
	}

	err = repo.Validate(ctx, "unknown-token")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Validate() error = %v, want %v", err, gorm.ErrRecordNotFound)
	}
}

func TestRegisterCreatesTokenHashIndex(t *testing.T) {
	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := Register(db); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasIndex(&Token{}, "idx_tokens_hash") {
		t.Fatal("token hash index was not created")
	}

	first := Token{ID: uuid.MustOrdered(), Name: "first", Hash: "duplicate-hash"}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first token: %v", err)
	}
	duplicate := Token{ID: uuid.MustOrdered(), Name: "duplicate", Hash: first.Hash}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate token hash was accepted")
	}
}

func testTokenRepository(t *testing.T) *TokenRepository {
	t.Helper()

	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := Register(db); err != nil {
		t.Fatal(err)
	}
	return NewRepo(db)
}
