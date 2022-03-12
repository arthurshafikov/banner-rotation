package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
)

type Banners struct {
	db    *sqlx.DB
	table string
}

func NewBanners(db *sqlx.DB) *Banners {
	return &Banners{
		db:    db,
		table: "banners",
	}
}

func (b *Banners) AddBanner(ctx context.Context, description string) error {
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1)", b.table)
	if err := b.db.QueryRowContext(ctx, query, description).Err(); err != nil {
		return err
	}

	return nil
}

func (b *Banners) DeleteBanner(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Err(); err != nil {
		return err
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
