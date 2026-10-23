package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/crypto"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

type Token struct {
	database.BaseModel

	Name string `json:"name"`
	Hash string `json:"-"`
}

type TokenRepository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) List(ctx context.Context) ([]Token, error) {
	tokens, err := gorm.G[Token](r.db).Order("id DESC").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("list auth tokens: %w", err)
	}
	return tokens, nil
}

func (r *TokenRepository) Create(ctx context.Context, name string) (string, Token, error) {
	id, _ := uuid.Ordered()

	by := make([]byte, 64)
	if _, err := rand.Read(by); err != nil {
		return "", Token{}, fmt.Errorf("generate secret: %w", err)
	}
	tokenStr := hex.EncodeToString(by)
	hash := crypto.Sum(tokenStr)

	token := &Token{
		ID:   id,
		Name: name,
		Hash: hash,
	}

	if err := gorm.G[Token](r.db).Create(ctx, token); err != nil {
		return "", Token{}, fmt.Errorf("store token: %w", err)
	}

	return tokenStr, *token, nil
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[Token](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}

func Register(db *gorm.DB) error {
	return db.AutoMigrate(
		&Token{},
	)
}
