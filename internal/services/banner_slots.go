package services

import (
	"context"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type BannerSlotService struct {
	repo repository.BannerSlots
}

func NewBannerSlotService(repo repository.BannerSlots) *BannerSlotService {
	return &BannerSlotService{
		repo: repo,
	}
}

func (bs *BannerSlotService) AssociateBannerToSlot(ctx context.Context, bannerId, slotId int64) error {
	return bs.repo.AddBannerSlot(ctx, bannerId, slotId)
}

func (bs *BannerSlotService) DissociateBannerFromSlot(ctx context.Context, bannerId, slotId int64) error {
	return bs.repo.DeleteBannerSlot(ctx, bannerId, slotId)
}

func (bs *BannerSlotService) GetByBannerAndSlotIds(ctx context.Context, bannerId, slotId int64) (*core.BannerSlot, error) {
	return bs.repo.GetByBannerAndSlotIds(ctx, bannerId, slotId)
}

func (bs *BannerSlotService) GetRandomBannerIdExceptExcluded(
	ctx context.Context,
	slotId,
	excludedBannerId int64,
) (int64, error) {
	return bs.repo.GetRandomBannerIdExceptExcluded(ctx, slotId, excludedBannerId)
}
