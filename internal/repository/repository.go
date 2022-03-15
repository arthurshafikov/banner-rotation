package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository/repositories"
)

type Banners interface {
	AddBanner(ctx context.Context, description string) error
	DeleteBanner(ctx context.Context, id int64) error
	GetBanner(ctx context.Context, id int64) (*core.Banner, error)
}

type Slots interface {
	AddSlot(ctx context.Context, description string) error
	DeleteSlot(ctx context.Context, id int64) error
	GetSlot(ctx context.Context, id int64) (*core.Slot, error)
}

type BannerSlots interface {
	AddBannerSlot(ctx context.Context, bannerId, slotId int64) error
	DeleteBannerSlot(ctx context.Context, bannerId, slotId int64) error
	GetByBannerAndSlotIds(ctx context.Context, bannerId, slotId int64) (*core.BannerSlot, error)
	GetRandomBannerIdExceptExcluded(ctx context.Context, slotId, excludedBannerId int64) (int64, error)
}

type SocialGroups interface {
	AddSocialGroup(ctx context.Context, description string) error
	DeleteSocialGroup(ctx context.Context, id int64) error
	GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error)
}

type BannerSlotSocialGroups interface {
	IncrementClick(ctx context.Context, bannerSlotId, socialGroupId int64) error
	IncrementView(ctx context.Context, bannerSlotId, socialGroupId int64) error
	GetTheMostProfitableBannerId(ctx context.Context, slotId, socialGroupId int64) (int64, error)
}

type Repository struct {
	Banners
	Slots
	BannerSlots
	SocialGroups
	BannerSlotSocialGroups
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Banners:                repositories.NewBanners(db),
		Slots:                  repositories.NewSlots(db),
		BannerSlots:            repositories.NewBannerSlots(db),
		SocialGroups:           repositories.NewSocialGroups(db),
		BannerSlotSocialGroups: repositories.NewBannerSlotSocialGroups(db),
	}
}
