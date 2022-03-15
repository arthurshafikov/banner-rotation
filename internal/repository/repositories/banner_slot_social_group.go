package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
)

type BannerSlotSocialGroups struct {
	db    *sqlx.DB
	table string
}

func NewBannerSlotSocialGroups(db *sqlx.DB) *BannerSlotSocialGroups {
	return &BannerSlotSocialGroups{
		db:    db,
		table: core.BannerSlotSocialGroupTable,
	}
}

func (bss *BannerSlotSocialGroups) IncrementClick(ctx context.Context, bannerSlotId, socialGroupId int64) error {
	bannerSlotSocialGroup := &core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId,
		SocialGroupId: socialGroupId,
	}
	err := bss.scanFirst(ctx, bannerSlotSocialGroup)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return err
	}

	if bannerSlotSocialGroup.ID == 0 {
		err = bss.createAndScan(ctx, bannerSlotSocialGroup)
		if err != nil {
			return err
		}
	}

	query := fmt.Sprintf("UPDATE %s SET clicks = clicks + 1 WHERE id=$1;", bss.table)
	if err := bss.db.QueryRowContext(ctx, query, bannerSlotSocialGroup.ID).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (bss *BannerSlotSocialGroups) IncrementView(ctx context.Context, bannerSlotId, socialGroupId int64) error {
	bannerSlotSocialGroup := &core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId,
		SocialGroupId: socialGroupId,
	}
	err := bss.scanFirst(ctx, bannerSlotSocialGroup)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return err
	}

	if bannerSlotSocialGroup.ID == 0 {
		err = bss.createAndScan(ctx, bannerSlotSocialGroup)
		if err != nil {
			return err
		}
	}

	query := fmt.Sprintf("UPDATE %s SET views = views + 1 WHERE id = $1;", bss.table)
	if err := bss.db.QueryRowContext(ctx, query, bannerSlotSocialGroup.ID).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (bss *BannerSlotSocialGroups) GetTheMostProfitableBannerId(
	ctx context.Context,
	slotId,
	socialGroupId int64,
) (int64, error) {
	var result = struct {
		Banner_id int64   `db:"banner_id"`
		Ctr       float64 `db:"ctr"`
	}{}
	query := fmt.Sprintf(
		`SELECT %[2]s.banner_id, CAST(%[1]s.clicks AS DECIMAL)/%[1]s.views as ctr FROM %[1]s
			LEFT JOIN %[2]s ON %[2]s.id = %[1]s.banner_slot_id
			WHERE %[2]s.slot_id = $1 AND %[1]s.social_group_id = $2
			ORDER BY ctr DESC
			LIMIT 1;`,
		bss.table,
		core.BannerSlotTable,
	)
	if err := bss.db.GetContext(ctx, &result, query, slotId, socialGroupId); err != nil {
		return 0, err
	}

	return result.Banner_id, nil
}

func (bss *BannerSlotSocialGroups) GetRandomExceptExcludedBannerId(
	ctx context.Context,
	slotId,
	excludedBannerId int64,
) (int64, error) {
	var result = struct {
		Banner_id int64 `db:"banner_id"`
	}{}
	// https://www.gab.lc/articles/bigdata_postgresql_order_by_random/
	query := fmt.Sprintf(
		`SELECT %[1]s.banner_id FROM %[1]s
			WHERE %[1]s.slot_id = $1 AND %[1]s.banner_id != $2
			ORDER BY RANDOM()
			LIMIT 1;`,
		core.BannerSlotTable,
	)
	if err := bss.db.GetContext(ctx, &result, query, slotId, excludedBannerId); err != nil {
		return 0, err
	}

	return result.Banner_id, nil
}

func (bss *BannerSlotSocialGroups) scanFirst(
	ctx context.Context,
	bannerSlotSocialGroup *core.BannerSlotSocialGroup,
) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE banner_slot_id=$1 AND social_group_id=$2 LIMIT 1;", bss.table)

	if err := bss.db.GetContext(
		ctx,
		bannerSlotSocialGroup,
		query,
		bannerSlotSocialGroup.BannerSlotId,
		bannerSlotSocialGroup.SocialGroupId,
	); err != nil {
		return err
	}

	return nil
}

func (bss *BannerSlotSocialGroups) createAndScan(
	ctx context.Context,
	bannerSlotSocialGroup *core.BannerSlotSocialGroup,
) error {
	query := fmt.Sprintf("INSERT INTO %s (banner_slot_id, social_group_id) VALUES ($1, $2) RETURNING id;", bss.table)
	err := bss.db.QueryRowxContext(
		ctx,
		query,
		bannerSlotSocialGroup.BannerSlotId,
		bannerSlotSocialGroup.SocialGroupId,
	).Scan(&bannerSlotSocialGroup.ID)
	if err != nil {
		return err
	}

	return nil
}
