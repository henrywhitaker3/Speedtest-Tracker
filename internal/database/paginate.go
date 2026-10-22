package database

import (
	"context"
	"fmt"
	"math"

	"gorm.io/gorm"
)

type Pagination[T any] struct {
	Limit      int   `json:"limit,omitempty"`
	Page       int   `json:"page,omitempty"`
	Total      int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Data       []T   `json:"data"`
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[T]) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 25
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func Paginate[T any](
	ctx context.Context,
	pagination *Pagination[T],
	db *gorm.DB,
) (func(db *gorm.Statement), error) {
	count, err := gorm.G[T](db).Count(ctx, "id")
	if err != nil {
		return nil, fmt.Errorf("get id count: %w", err)
	}

	pagination.Total = count
	pagination.TotalPages = int(math.Ceil(float64(count) / float64(pagination.GetLimit())))

	return func(db *gorm.Statement) {
		db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}, nil
}
