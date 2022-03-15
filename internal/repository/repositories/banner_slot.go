package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
)

type BannerSlots struct {
	db    *sqlx.DB
	table string
}

func NewBannerSlots(db *sqlx.DB) *BannerSlots {
	return &BannerSlots{
		db:    db,
		table: core.BannerSlotTable,
	}
}

func (bs *BannerSlots) AddBannerSlot(ctx context.Context, bannerId, slotId int64) error {
	query := fmt.Sprintf("INSERT INTO %s (banner_id, slot_id) VALUES ($1, $2)", bs.table)
	if err := bs.db.QueryRowContext(ctx, query, bannerId, slotId).Err(); err != nil {
		return err
	}

	return nil
}

func (bs *BannerSlots) DeleteBannerSlot(ctx context.Context, bannerId, slotId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE banner_id=$1 AND slot_id=$2", bs.table)
	if err := bs.db.QueryRowContext(ctx, query, bannerId, slotId).Err(); err != nil {
		return err
	}

	return nil
}

func (bs *BannerSlots) GetByServiceAndBannerIds(ctx context.Context, bannerId, slotId int64) (*core.BannerSlot, error) {
	bannerSlot := core.BannerSlot{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE banner_id=$1 AND slot_id=$2", bs.table)
	if err := bs.db.GetContext(ctx, &bannerSlot, query, bannerId, slotId); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return &bannerSlot, nil
}
