// Package results
package results

import (
	"context"
	"fmt"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"gorm.io/gorm"
)

type Result struct {
	database.BaseModel

	Ping       *float64 `json:"ping"`
	Download   *int64   `json:"download"`
	Upload     *int64   `json:"upload"`
	ServerID   *int     `json:"server_id"`
	ServerHost *string  `json:"server_host"`
	ServerName *string  `json:"server_name"`
	URL        *string  `json:"url"`
	Error      *string  `json:"error,omitempty"`
	Complete   bool     `json:"complete"        gorm:"index:idx_results_timeseries,expression:created_at,where:complete = true AND deleted_at IS NULL"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Latest(ctx context.Context) (*Result, error) {
	result, err := gorm.G[Result](r.db).Order("id DESC").First(ctx)
	return &result, err
}

func (r *Repository) Create(ctx context.Context) (*Result, error) {
	result := &Result{
		ID: uuid.MustOrdered(),
	}
	if err := gorm.G[Result](r.db).Create(ctx, result); err != nil {
		return nil, fmt.Errorf("store initial result: %w", err)
	}

	return result, nil
}

func (r *Repository) Fresh(ctx context.Context, result *Result) error {
	fresh, err := r.Get(ctx, result.ID)
	if err != nil {
		return err
	}
	*result = *fresh
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Result, error) {
	result, err := gorm.G[Result](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("get result: %w", err)
	}
	return &result, nil
}

func (r *Repository) Update(ctx context.Context, result *Result) error {
	if _, err := gorm.G[Result](r.db).
		Where("id = ?", result.ID).
		Select("*").
		Updates(ctx, *result); err != nil {
		return fmt.Errorf("store updated result: %w", err)
	}
	return nil
}

type ListOpts struct {
	Page  int
	Limit int
}

func (r *Repository) List(
	ctx context.Context,
	opts ListOpts,
) (*database.Pagination[Result], error) {
	out := &database.Pagination[Result]{
		Page:  opts.Page,
		Limit: opts.Limit,
	}
	paginator, err := database.Paginate(ctx, out, r.db)
	if err != nil {
		return nil, fmt.Errorf("build results paginator: %w", err)
	}

	results, err := gorm.G[Result](r.db).Scopes(paginator).Order("id DESC").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("list results paginated: %w", err)
	}
	out.Data = results

	return out, nil
}

type Stat struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Ping      *float64  `json:"ping"`
	Download  *int64    `json:"download"`
	Upload    *int64    `json:"upload"`
}

func (r *Repository) Timeseries(ctx context.Context, period time.Duration) ([]Stat, error) {
	stats, err := gorm.G[Result](r.db).
		Select("id, ping, download, upload, created_at").
		Where("created_at >= ?", time.Now().Add(-period)).
		Where("created_at <= ?", time.Now()).
		Where("complete = true").
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("get result timeseries: %w", err)
	}

	out := []Stat{}
	for _, s := range stats {
		out = append(out, Stat{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			Ping:      s.Ping,
			Download:  s.Download,
			Upload:    s.Upload,
		})
	}

	return out, nil
}

func Register(db *gorm.DB) error {
	return db.AutoMigrate(
		&Result{},
	)
}
