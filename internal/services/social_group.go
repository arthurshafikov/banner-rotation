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

func (sg *SocialGroupService) AddSocialGroup(ctx context.Context, description string) (int64, error) {
	return sg.repo.AddSocialGroup(ctx, description)
}

func (sg *SocialGroupService) DeleteSocialGroup(ctx context.Context, id int64) error {
	return sg.repo.DeleteSocialGroup(ctx, id)
}

func (sg *SocialGroupService) GetSocialGroup(ctx context.Context, id int64) (*core.SocialGroup, error) {
	return sg.repo.GetSocialGroup(ctx, id)
}
