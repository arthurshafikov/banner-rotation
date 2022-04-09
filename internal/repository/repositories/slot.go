package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/jmoiron/sqlx"
)

type Slots struct {
	db    *sqlx.DB
	table string
}

func NewSlots(db *sqlx.DB) *Slots {
	return &Slots{
		db:    db,
		table: core.SlotsTable,
	}
}

func (b *Slots) AddSlot(ctx context.Context, description string) (int64, error) {
	var slotId int64
	query := fmt.Sprintf("INSERT INTO %s (description) VALUES ($1) RETURNING id;", b.table)
	err := b.db.QueryRowxContext(ctx, query, description).Scan(&slotId)

	return slotId, err
}

func (b *Slots) DeleteSlot(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", b.table)
	if err := b.db.QueryRowContext(ctx, query, id).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
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
