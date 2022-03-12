package services

import "github.com/thewolf27/banner-rotation/internal/repository"

type BannerSlotSocialGroupService struct {
	repo repository.BannerSlotSocialGroups
}

func NewBannerSlotSocialGroupService(repo repository.BannerSlotSocialGroups) *BannerSlotSocialGroupService {
	return &BannerSlotSocialGroupService{
		repo: repo,
	}
}
