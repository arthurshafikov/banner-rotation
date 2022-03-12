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

type BannerSlotSocialGroups interface {
}

type Repository struct {
	Banners
	Slots
	BannerSlotSocialGroups
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Banners:                repositories.NewBanners(db),
		Slots:                  repositories.NewSlots(db),
		BannerSlotSocialGroups: repositories.NewBannerSlotSocialGroups(db),
	}
}
