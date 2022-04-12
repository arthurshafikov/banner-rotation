package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/jmoiron/sqlx"
)

type SocialGroups struct {
	db    *sqlx.DB
	table string
}

func NewSocialGroups(db *sqlx.DB) *SocialGroups {
	return &SocialGroups{
		db:    db,
		table: core.SocialGroupTable,
	}
}

func (b *SocialGroups) AddSocialGroup(ctx context.Context, description string) (int64, error) {
	var socialGroupID int64
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1) RETURNING id;", b.table)
	err := b.db.QueryRowxContext(ctx, query, description).Scan(&socialGroupID)

	return socialGroupID, err
}

func (b *SocialGroups) DeleteSocialGroup(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
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
