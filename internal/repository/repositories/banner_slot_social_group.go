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

func (bssg *BannerSlotSocialGroups) IncrementClick(ctx context.Context, bannerSlotId, socialGroupId int64) error {
	bannerSlotSocialGroup, err := bssg.firstOrCreate(ctx, bannerSlotId, socialGroupId)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET clicks = clicks + 1 WHERE id = $1;", bssg.table)
	if err := bssg.db.QueryRowContext(ctx, query, bannerSlotSocialGroup.ID).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (bssg *BannerSlotSocialGroups) IncrementView(ctx context.Context, bannerSlotId, socialGroupId int64) error {
	bannerSlotSocialGroup, err := bssg.firstOrCreate(ctx, bannerSlotId, socialGroupId)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET views = views + 1 WHERE id = $1;", bssg.table)
	if err := bssg.db.QueryRowContext(ctx, query, bannerSlotSocialGroup.ID).Scan(); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return err
		}
	}

	return nil
}

func (bssg *BannerSlotSocialGroups) GetTheMostProfitableBannerId(
	ctx context.Context,
	slotId,
	socialGroupId int64,
) (int64, error) {
	var result = struct {
		BannerId int64   `db:"banner_id"`
		Ctr      float64 `db:"ctr"`
	}{}
	escapeDivideByZero := fmt.Sprintf(`
		CASE %[1]s.views
			WHEN 0 THEN 1
			ELSE %[1]s.views
		END
	`, bssg.table)
	query := fmt.Sprintf(
		`SELECT %[2]s.banner_id, CAST(%[1]s.clicks AS DECIMAL)/%[3]s as ctr FROM %[1]s
			LEFT JOIN %[2]s ON %[2]s.id = %[1]s.banner_slot_id
			WHERE %[2]s.slot_id = $1 AND %[1]s.social_group_id = $2
			ORDER BY ctr DESC
			LIMIT 1;`,
		bssg.table,
		core.BannerSlotTable,
		escapeDivideByZero,
	)

	if err := bssg.db.GetContext(ctx, &result, query, slotId, socialGroupId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	return result.BannerId, nil
}

func (bssg *BannerSlotSocialGroups) firstOrCreate(
	ctx context.Context,
	bannerSlotId,
	socialGroupId int64,
) (*core.BannerSlotSocialGroup, error) {
	bannerSlotSocialGroup := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId,
		SocialGroupId: socialGroupId,
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE banner_slot_id=$1 AND social_group_id=$2 LIMIT 1;", bssg.table)
	if err := bssg.db.GetContext(
		ctx,
		&bannerSlotSocialGroup.ID,
		query,
		bannerSlotId,
		socialGroupId,
	); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, err
	} else if errors.Is(sql.ErrNoRows, err) {
		query := fmt.Sprintf("INSERT INTO %s (banner_slot_id, social_group_id) VALUES ($1, $2) RETURNING id;", bssg.table)

		if err := bssg.db.QueryRowxContext(
			ctx,
			query,
			bannerSlotId,
			socialGroupId,
		).Scan(&bannerSlotSocialGroup.ID); err != nil {
			return nil, err
		}
	}

	return &bannerSlotSocialGroup, nil
}
