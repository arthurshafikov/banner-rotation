package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
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
	bannerSlot, err := bssg.bannerSlotService.GetByBannerAndSlotIDs(ctx, inp.BannerID, inp.SlotID)
	if err != nil {
		return err
	}

	if err := bssg.repo.IncrementClick(ctx, bannerSlot.ID, inp.SocialGroupID); err != nil {
		return err
	}

	return bssg.queue.AddToQueue("clicks", core.IncrementEvent{
		BannerID:      inp.BannerID,
		SlotID:        inp.SlotID,
		SocialGroupID: inp.SocialGroupID,
		Datetime:      time.Now(),
	})
}

func (bssg *BannerSlotSocialGroupService) GetBannerIDToShow(
	ctx context.Context,
	inp core.GetBannerRequest,
) (int64, error) {
	bannerID, err := bssg.repo.GetTheMostProfitableBannerID(ctx, inp.SlotID, inp.SocialGroupID)
	if err != nil {
		return 0, err
	}

	if bannerID == 0 || bssg.rollADice(bssg.eGreedValue) {
		bannerID, err = bssg.bannerSlotService.GetRandomBannerIDExceptExcluded(ctx, inp.SlotID, bannerID)
		if err != nil {
			return 0, err
		}
	}

	bannerSlot, err := bssg.bannerSlotService.GetByBannerAndSlotIDs(ctx, bannerID, inp.SlotID)
	if err != nil {
		return 0, err
	}
	if err := bssg.repo.IncrementView(ctx, bannerSlot.ID, inp.SocialGroupID); err != nil {
		return 0, err
	}

	if err := bssg.queue.AddToQueue("views", core.IncrementEvent{
		BannerID:      bannerID,
		SlotID:        inp.SlotID,
		SocialGroupID: inp.SocialGroupID,
		Datetime:      time.Now(),
	}); err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (bssg *BannerSlotSocialGroupService) rollADice(chance float64) bool {
	return rand.Float64() <= chance //nolint
}
