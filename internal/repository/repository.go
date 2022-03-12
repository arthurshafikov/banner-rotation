package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/thewolf27/banner-rotation/internal/repository/repositories"
)

type Banners interface {
}

type BannerSlotSocialGroups interface {
}

type Repository struct {
	Banners
	BannerSlotSocialGroups
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Banners:                repositories.NewBanners(db),
		BannerSlotSocialGroups: repositories.NewBannerSlotSocialGroups(db),
	}
}
