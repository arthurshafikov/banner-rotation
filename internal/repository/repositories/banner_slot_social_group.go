package repositories

import "github.com/jmoiron/sqlx"

type BannerSlotSocialGroups struct {
	db *sqlx.DB
}

func NewBannerSlotSocialGroups(db *sqlx.DB) *BannerSlotSocialGroups {
	return &BannerSlotSocialGroups{
		db: db,
	}
}
