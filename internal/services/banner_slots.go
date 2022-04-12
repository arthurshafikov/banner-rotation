package services

import (
	"context"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
)

type BannerSlotService struct {
	repo repository.BannerSlots
}

func NewBannerSlotService(repo repository.BannerSlots) *BannerSlotService {
	return &BannerSlotService{
		repo: repo,
	}
}

func (bs *BannerSlotService) AssociateBannerToSlot(ctx context.Context, bannerID, slotID int64) (int64, error) {
	return bs.repo.AddBannerSlot(ctx, bannerID, slotID)
}

func (bs *BannerSlotService) DissociateBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	return bs.repo.DeleteBannerSlot(ctx, bannerID, slotID)
}

func (bs *BannerSlotService) GetByBannerAndSlotIDs(
	ctx context.Context,
	bannerID,
	slotID int64,
) (*core.BannerSlot, error) {
	return bs.repo.GetByBannerAndSlotIDs(ctx, bannerID, slotID)
}

func (bs *BannerSlotService) GetRandomBannerIDExceptExcluded(
	ctx context.Context,
	slotID,
	excludedBannerID int64,
) (int64, error) {
	return bs.repo.GetRandomBannerIDExceptExcluded(ctx, slotID, excludedBannerID)
}
