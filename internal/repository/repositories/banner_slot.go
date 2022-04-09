package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/jmoiron/sqlx"
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

func (bs *BannerSlots) AddBannerSlot(ctx context.Context, bannerId, slotId int64) (int64, error) {
	var bannerSlotId int64
	query := fmt.Sprintf("INSERT INTO %s (banner_id, slot_id) VALUES ($1, $2) RETURNING id;", bs.table)
	err := bs.db.QueryRowxContext(ctx, query, bannerId, slotId).Scan(&bannerSlotId)

	return bannerSlotId, err
}

func (bs *BannerSlots) DeleteBannerSlot(ctx context.Context, bannerId, slotId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE banner_id=$1 AND slot_id=$2", bs.table)
	if err := bs.db.QueryRowContext(ctx, query, bannerId, slotId).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (bs *BannerSlots) GetByBannerAndSlotIds(ctx context.Context, bannerId, slotId int64) (*core.BannerSlot, error) {
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

func (bs *BannerSlots) GetRandomBannerIdExceptExcluded(
	ctx context.Context,
	slotId,
	excludedBannerId int64,
) (int64, error) {
	var result = struct {
		BannerId int64 `db:"banner_id"`
	}{}
	query := fmt.Sprintf(
		`SELECT banner_id FROM %s
			WHERE slot_id = $1 AND banner_id != $2
			ORDER BY RANDOM()
			LIMIT 1;`,
		bs.table,
	)
	if err := bs.db.GetContext(ctx, &result, query, slotId, excludedBannerId); err != nil {
		return 0, err
	}

	return result.BannerId, nil
}
