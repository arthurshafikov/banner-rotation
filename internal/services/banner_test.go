package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_repository "github.com/arthurshafikov/banner-rotation/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getBannerRepoMock(t *testing.T) (context.Context, *mock_repository.MockBanners) {
	t.Helper()
	ctl := gomock.NewController(t)
	bannerRepo := mock_repository.NewMockBanners(ctl)

	return context.Background(), bannerRepo
}

func TestAddBanner(t *testing.T) {
	ctx, bannerRepo := getBannerRepoMock(t)
	expectedBanner := core.Banner{
		ID: 5,
	}
	gomock.InOrder(
		bannerRepo.EXPECT().AddBanner(ctx, "test_description").Return(expectedBanner.ID, nil),
	)
	b := NewBannerService(bannerRepo)

	bannerId, err := b.AddBanner(ctx, "test_description")

	require.NoError(t, err)
	require.Equal(t, expectedBanner.ID, bannerId)
}

func TestDeleteBanner(t *testing.T) {
	ctx, bannerRepo := getBannerRepoMock(t)
	gomock.InOrder(
		bannerRepo.EXPECT().DeleteBanner(ctx, int64(23)).Return(nil),
	)
	b := NewBannerService(bannerRepo)

	err := b.DeleteBanner(ctx, 23)

	require.NoError(t, err)
}

func TestGetBanner(t *testing.T) {
	ctx, bannerRepo := getBannerRepoMock(t)
	expected := &core.Banner{
		ID: 22,
	}
	gomock.InOrder(
		bannerRepo.EXPECT().GetBanner(ctx, expected.ID).Return(expected, nil),
	)
	b := NewBannerService(bannerRepo)

	result, err := b.GetBanner(ctx, 22)

	require.Equal(t, expected, result)
	require.NoError(t, err)
}
