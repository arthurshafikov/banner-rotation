package services

import (
	"context"
	"math/rand"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type BannerSlotSocialGroupService struct {
	repo repository.BannerSlotSocialGroups

	bannerSlotService *BannerSlotService
}

func NewBannerSlotSocialGroupService(
	repo repository.BannerSlotSocialGroups,
	bannerSlotService *BannerSlotService,
) *BannerSlotSocialGroupService {
	return &BannerSlotSocialGroupService{
		repo:              repo,
		bannerSlotService: bannerSlotService,
	}
}

func (bss *BannerSlotSocialGroupService) IncrementClick(ctx context.Context, inp core.IncrementClickInput) error {
	bannerSlot, err := bss.bannerSlotService.GetByServiceAndBannerIds(ctx, inp.BannerId, inp.SlotId)
	if err != nil {
		return err
	}

	return bss.repo.IncrementClick(ctx, bannerSlot.ID, inp.SocialGroupId)
}

func (bss *BannerSlotSocialGroupService) GetBanner(ctx context.Context, inp core.GetBannerRequest) (int64, error) {
	bannerId, err := bss.repo.GetTheMostProfitableBannerId(ctx, inp.SlotId, inp.SocialGroupId)
	if err != nil {
		return 0, err
	}

	// todo config E-greed value
	if bss.rollADice(0.1) {
		bannerId, err = bss.repo.GetRandomExceptExcludedBannerId(ctx, inp.SlotId, bannerId)
		if err != nil {
			return 0, err
		}
	}

	bannerSlot, err := bss.bannerSlotService.GetByServiceAndBannerIds(ctx, bannerId, inp.SlotId)
	if err != nil {
		return 0, err
	}
	if err := bss.repo.IncrementView(ctx, bannerSlot.ID, inp.SocialGroupId); err != nil {
		return 0, err
	}

	return bannerId, nil
}

func (bss *BannerSlotSocialGroupService) rollADice(chance float64) bool {
	return rand.Float64() <= chance
}
