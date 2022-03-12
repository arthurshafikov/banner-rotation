package services

import (
	"context"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type BannerService struct {
	repo repository.Banners
}

func NewBannerService(repo repository.Banners) *BannerService {
	return &BannerService{
		repo: repo,
	}
}

func (b *BannerService) AddBanner(ctx context.Context, description string) error {
	return b.repo.AddBanner(ctx, description)
}

func (b *BannerService) DeleteBanner(ctx context.Context, id int64) error {
	return b.repo.DeleteBanner(ctx, id)
}

func (b *BannerService) GetBanner(ctx context.Context, id int64) (*core.Banner, error) {
	return b.repo.GetBanner(ctx, id)
}
