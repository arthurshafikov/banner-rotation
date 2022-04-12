package services

import (
	"context"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
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
	AssociateBannerToSlot(ctx context.Context, bannerID, slotID int64) (int64, error)
	DissociateBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	GetByBannerAndSlotIDs(ctx context.Context, bannerID, slotID int64) (*core.BannerSlot, error)
	GetRandomBannerIDExceptExcluded(ctx context.Context, slotID, excludedBannerID int64) (int64, error)
}

type SocialGroups interface {
	AddSocialGroup(ctx context.Context, description string) (int64, error)
	DeleteSocialGroup(ctx context.Context, id int64) error
	GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error)
}

type BannerSlotSocialGroups interface {
	IncrementClick(ctx context.Context, inp core.IncrementClickInput) error
	GetBannerIDToShow(ctx context.Context, inp core.GetBannerRequest) (int64, error)
}

type Services struct {
	Banners
	Slots
	BannerSlots
	SocialGroups
	BannerSlotSocialGroups
}

type Queue interface {
	AddToQueue(topic string, value interface{}) error
}

type Dependencies struct {
	Repository  *repository.Repository
	EGreedValue float64
	Queue
}

func NewServices(deps Dependencies) *Services {
	bannerService := NewBannerService(deps.Repository.Banners)
	slotService := NewSlotService(deps.Repository.Slots)
	bannerSlotService := NewBannerSlotService(deps.Repository.BannerSlots)
	socialGroupService := NewSocialGroupService(deps.Repository.SocialGroups)
	bannerSlotSocialGroupService := NewBannerSlotSocialGroupService(
		deps.Repository.BannerSlotSocialGroups,
		bannerSlotService,
		deps.EGreedValue,
		deps.Queue,
	)

	return &Services{
		Banners:                bannerService,
		Slots:                  slotService,
		BannerSlots:            bannerSlotService,
		SocialGroups:           socialGroupService,
		BannerSlotSocialGroups: bannerSlotSocialGroupService,
	}
}
