package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type BannerSlotSocialGroupService struct {
	repo        repository.BannerSlotSocialGroups
	eGreedValue float64

	bannerSlotService BannerSlots
	queue             Queue
}

func NewBannerSlotSocialGroupService(
	repo repository.BannerSlotSocialGroups,
	bannerSlotService BannerSlots,
	eGreedValue float64,
	queue Queue,
) *BannerSlotSocialGroupService {
	return &BannerSlotSocialGroupService{
		repo:              repo,
		eGreedValue:       eGreedValue,
		bannerSlotService: bannerSlotService,
		queue:             queue,
	}
}

func (bssg *BannerSlotSocialGroupService) IncrementClick(ctx context.Context, inp core.IncrementClickInput) error {
	bannerSlot, err := bssg.bannerSlotService.GetByBannerAndSlotIds(ctx, inp.BannerId, inp.SlotId)
	if err != nil {
		return err
	}

	if err := bssg.repo.IncrementClick(ctx, bannerSlot.ID, inp.SocialGroupId); err != nil {
		return nil
	}

	return bssg.queue.AddToQueue("clicks", core.IncrementEvent{
		BannerId:      inp.BannerId,
		SlotId:        inp.SlotId,
		SocialGroupId: inp.SocialGroupId,
		Datetime:      time.Now(),
	})
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

	if err := bssg.queue.AddToQueue("views", core.IncrementEvent{
		BannerId:      bannerId,
		SlotId:        inp.SlotId,
		SocialGroupId: inp.SocialGroupId,
		Datetime:      time.Now(),
	}); err != nil {
		return 0, err
	}

	return bannerId, nil
}

func (bssg *BannerSlotSocialGroupService) rollADice(chance float64) bool {
	return rand.Float64() <= chance
}
