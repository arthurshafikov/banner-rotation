package services

import (
	"context"

	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
)

type SocialGroupService struct {
	repo repository.SocialGroups
}

func NewSocialGroupService(repo repository.SocialGroups) *SocialGroupService {
	return &SocialGroupService{
		repo: repo,
	}
}

func (b *SocialGroupService) AddSocialGroup(ctx context.Context, description string) error {
	return b.repo.AddSocialGroup(ctx, description)
}

func (b *SocialGroupService) DeleteSocialGroup(ctx context.Context, id int64) error {
	return b.repo.DeleteSocialGroup(ctx, id)
}

func (b *SocialGroupService) GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error) {
	return b.repo.GetSocialGroup(ctx, id)
}
