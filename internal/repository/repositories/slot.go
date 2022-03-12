package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
)

type Slots struct {
	db    *sqlx.DB
	table string
}

func NewSlots(db *sqlx.DB) *Slots {
	return &Slots{
		db:    db,
		table: "slots",
	}
}

func (b *Slots) AddSlot(ctx context.Context, description string) error {
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1)", b.table)
	if err := b.db.QueryRowContext(ctx, query, description).Err(); err != nil {
		return err
	}

	return nil
}

func (b *Slots) DeleteSlot(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Err(); err != nil {
		return err
	}

	return nil
}

func (b *Slots) GetSlot(ctx context.Context, id int64) (*core.Slot, error) {
	slot := core.Slot{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", b.table)
	if err := b.db.GetContext(ctx, &slot, query, id); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return &slot, nil
}
