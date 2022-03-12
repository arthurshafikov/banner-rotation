package services

import (
	"context"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
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
	AssociateBannerToSlot(ctx context.Context, bannerId, slotId int64) error
	DissociateBannerFromSlot(ctx context.Context, bannerId, slotId int64) error
}

type BannerSlotSocialGroups interface {
}

type Services struct {
	Banners
	Slots
	BannerSlots
	BannerSlotSocialGroups
}

type Dependencies struct {
	Repository *repository.Repository
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Banners:                NewBannerService(deps.Repository.Banners),
		Slots:                  NewSlotService(deps.Repository.Slots),
		BannerSlots:            NewBannerSlotService(deps.Repository.BannerSlots),
		BannerSlotSocialGroups: NewBannerSlotSocialGroupService(deps.Repository.BannerSlotSocialGroups),
	}
}
