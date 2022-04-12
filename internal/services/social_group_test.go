package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_repository "github.com/arthurshafikov/banner-rotation/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getSocialGroupRepoMock(t *testing.T) (context.Context, *mock_repository.MockSocialGroups) {
	t.Helper()
	ctl := gomock.NewController(t)
	socialGroupRepo := mock_repository.NewMockSocialGroups(ctl)

	return context.Background(), socialGroupRepo
}

func TestAddSocialGroup(t *testing.T) {
	ctx, socialGroupRepo := getSocialGroupRepoMock(t)
	expectedSocialGroup := core.SocialGroup{
		ID: 5,
	}
	gomock.InOrder(
		socialGroupRepo.EXPECT().AddSocialGroup(ctx, "test_description").Return(expectedSocialGroup.ID, nil),
	)
	sg := NewSocialGroupService(socialGroupRepo)

	socialGroupID, err := sg.AddSocialGroup(ctx, "test_description")

	require.NoError(t, err)
	require.Equal(t, expectedSocialGroup.ID, socialGroupID)
}

func TestDeleteSocialGroup(t *testing.T) {
	ctx, socialGroupRepo := getSocialGroupRepoMock(t)
	gomock.InOrder(
		socialGroupRepo.EXPECT().DeleteSocialGroup(ctx, int64(23)).Return(nil),
	)
	sg := NewSocialGroupService(socialGroupRepo)

	err := sg.DeleteSocialGroup(ctx, 23)

	require.NoError(t, err)
}

func TestGetSocialGroup(t *testing.T) {
	ctx, socialGroupRepo := getSocialGroupRepoMock(t)
	expected := &core.SocialGroup{
		ID: 22,
	}
	gomock.InOrder(
		socialGroupRepo.EXPECT().GetSocialGroup(ctx, expected.ID).Return(expected, nil),
	)
	sg := NewSocialGroupService(socialGroupRepo)

	result, err := sg.GetSocialGroup(ctx, 22)

	require.Equal(t, expected, result)
	require.NoError(t, err)
}
