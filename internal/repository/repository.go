package repository

import (
	"context"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/repository/repositories"
	"github.com/jmoiron/sqlx"
)

type Banners interface {
	AddBanner(ctx context.Context, description string) (int64, error)
	DeleteBanner(ctx context.Context, id int64) error
	GetBanner(ctx context.Context, id int64) (*core.Banner, error)
}

type Slots interface {
	AddSlot(ctx context.Context, description string) (int64, error)
	DeleteSlot(ctx context.Context, id int64) error
	GetSlot(ctx context.Context, id int64) (*core.Slot, error)
}

type BannerSlots interface {
	AddBannerSlot(ctx context.Context, bannerID, slotID int64) (int64, error)
	DeleteBannerSlot(ctx context.Context, bannerID, slotID int64) error
	GetByBannerAndSlotIDs(ctx context.Context, bannerID, slotID int64) (*core.BannerSlot, error)
	GetRandomBannerIDExceptExcluded(ctx context.Context, slotID, excludedBannerID int64) (int64, error)
}

type SocialGroups interface {
	AddSocialGroup(ctx context.Context, description string) (int64, error)
	DeleteSocialGroup(ctx context.Context, id int64) error
	GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error)
}

type BannerSlotSocialGroups interface {
	IncrementClick(ctx context.Context, bannerSlotID, socialGroupID int64) error
	IncrementView(ctx context.Context, bannerSlotID, socialGroupID int64) error
	GetTheMostProfitableBannerID(ctx context.Context, slotID, socialGroupID int64) (int64, error)
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
