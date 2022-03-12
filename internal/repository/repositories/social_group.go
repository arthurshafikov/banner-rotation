package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
)

type SocialGroups struct {
	db    *sqlx.DB
	table string
}

func NewSocialGroups(db *sqlx.DB) *SocialGroups {
	return &SocialGroups{
		db:    db,
		table: "social_groups",
	}
}

func (b *SocialGroups) AddSocialGroup(ctx context.Context, description string) error {
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1)", b.table)
	if err := b.db.QueryRowContext(ctx, query, description).Err(); err != nil {
		return err
	}

	return nil
}

func (b *SocialGroups) DeleteSocialGroup(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Err(); err != nil {
		return err
	}

	return nil
}

func (b *SocialGroups) GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error) {
	socialGroup := core.SocialGroup{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", b.table)
	if err := b.db.GetContext(ctx, &socialGroup, query, id); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return &socialGroup, nil
}
