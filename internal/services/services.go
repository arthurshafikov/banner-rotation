package services

import "github.com/thewolf27/banner-rotation/internal/repository"

type Banners interface {
}

type BannerSlotSocialGroups interface {
}

type Services struct {
	Banners                Banners
	BannerSlotSocialGroups BannerSlotSocialGroups
}

type Dependencies struct {
	Repository *repository.Repository
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Banners:                NewBannerService(deps.Repository.Banners),
		BannerSlotSocialGroups: NewBannerSlotSocialGroupService(deps.Repository.BannerSlotSocialGroups),
	}
}
