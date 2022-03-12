package services

import (
	"context"

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
