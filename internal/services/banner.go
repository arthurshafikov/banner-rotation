package services

import (
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
