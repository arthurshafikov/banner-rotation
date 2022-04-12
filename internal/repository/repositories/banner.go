package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/jmoiron/sqlx"
)

type Banners struct {
	db    *sqlx.DB
	table string
}

func NewBanners(db *sqlx.DB) *Banners {
	return &Banners{
		db:    db,
		table: core.BannersTable,
	}
}

func (b *Banners) AddBanner(ctx context.Context, description string) (int64, error) {
	var bannerID int64
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1) RETURNING id;", b.table)
	err := b.db.QueryRowxContext(ctx, query, description).Scan(&bannerID)

	return bannerID, err
}

func (b *Banners) DeleteBanner(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (b *Banners) GetBanner(ctx context.Context, id int64) (*core.Banner, error) {
	banner := core.Banner{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", b.table)
	if err := b.db.GetContext(ctx, &banner, query, id); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return &banner, nil
}
