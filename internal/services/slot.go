package services

import (
	"context"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
)

type SlotService struct {
	repo repository.Slots
}

func NewSlotService(repo repository.Slots) *SlotService {
	return &SlotService{
		repo: repo,
	}
}

func (b *SlotService) AddSlot(ctx context.Context, description string) (int64, error) {
	return b.repo.AddSlot(ctx, description)
}

func (b *SlotService) DeleteSlot(ctx context.Context, id int64) error {
	return b.repo.DeleteSlot(ctx, id)
}

func (b *SlotService) GetSlot(ctx context.Context, id int64) (*core.Slot, error) {
	return b.repo.GetSlot(ctx, id)
}
