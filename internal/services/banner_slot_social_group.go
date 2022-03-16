package services

import (
	"context"
	"math/rand"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type BannerSlotSocialGroupService struct {
	repo        repository.BannerSlotSocialGroups
	eGreedValue float64

	bannerSlotService BannerSlots
}

func NewBannerSlotSocialGroupService(
	repo repository.BannerSlotSocialGroups,
	bannerSlotService BannerSlots,
	eGreedValue float64,
) *BannerSlotSocialGroupService {
	return &BannerSlotSocialGroupService{
		repo:              repo,
		eGreedValue:       eGreedValue,
		bannerSlotService: bannerSlotService,
	}
}

func (bssg *BannerSlotSocialGroupService) IncrementClick(ctx context.Context, inp core.IncrementClickInput) error {
	bannerSlot, err := bssg.bannerSlotService.GetByBannerAndSlotIds(ctx, inp.BannerId, inp.SlotId)
	if err != nil {
		return err
	}

	return bssg.repo.IncrementClick(ctx, bannerSlot.ID, inp.SocialGroupId)
}

func (bssg *BannerSlotSocialGroupService) GetBannerIdToShow(
	ctx context.Context,
	inp core.GetBannerRequest,
) (int64, error) {
	bannerId, err := bssg.repo.GetTheMostProfitableBannerId(ctx, inp.SlotId, inp.SocialGroupId)
	if err != nil {
		return 0, err
	}

	if bannerId == 0 || bssg.rollADice(bssg.eGreedValue) {
		bannerId, err = bssg.bannerSlotService.GetRandomBannerIdExceptExcluded(ctx, inp.SlotId, bannerId)
		if err != nil {
			return 0, err
		}
	}

	bannerSlot, err := bssg.bannerSlotService.GetByBannerAndSlotIds(ctx, bannerId, inp.SlotId)
	if err != nil {
		return 0, err
	}
	if err := bssg.repo.IncrementView(ctx, bannerSlot.ID, inp.SocialGroupId); err != nil {
		return 0, err
	}

	return bannerId, nil
}

func (bssg *BannerSlotSocialGroupService) rollADice(chance float64) bool {
	return rand.Float64() <= chance
}
