package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
	mock_repository "github.com/thewolf27/banner-rotation/internal/repository/mocks"
)

func getBannerSlotRepoMock(t *testing.T) (context.Context, *mock_repository.MockBannerSlots) {
	t.Helper()
	ctl := gomock.NewController(t)
	bannerSlotRepo := mock_repository.NewMockBannerSlots(ctl)

	return context.Background(), bannerSlotRepo
}

func TestAssociateBannerToSlot(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	gomock.InOrder(
		bannerSlotRepo.EXPECT().AddBannerSlot(ctx, int64(1), int64(2)).Return(nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	err := bs.AssociateBannerToSlot(ctx, 1, 2)

	require.NoError(t, err)
}

func TestDissociateBannerFromSlot(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	gomock.InOrder(
		bannerSlotRepo.EXPECT().DeleteBannerSlot(ctx, int64(1), int64(2)).Return(nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	err := bs.DissociateBannerFromSlot(ctx, 1, 2)

	require.NoError(t, err)
}

func TestGetByBannerAndSlotIds(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	expected := &core.BannerSlot{
		ID:       5,
		BannerId: 1,
		SlotId:   2,
	}
	gomock.InOrder(
		bannerSlotRepo.EXPECT().GetByBannerAndSlotIds(ctx, expected.BannerId, expected.SlotId).
			Return(expected, nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	result, err := bs.GetByBannerAndSlotIds(ctx, 1, 2)

	require.Equal(t, expected, result)
	require.NoError(t, err)
}

func TestGetRandomBannerIdExceptExcluded(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	expected := &core.BannerSlot{
		ID:       5,
		BannerId: 1,
		SlotId:   3,
	}
	excludedBanner := core.Banner{
		ID: 56,
	}
	gomock.InOrder(
		bannerSlotRepo.EXPECT().
			GetRandomBannerIdExceptExcluded(ctx, expected.SlotId, excludedBanner.ID).
			Return(expected.ID, nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	result, err := bs.GetRandomBannerIdExceptExcluded(ctx, 3, 56)

	require.Equal(t, expected.ID, result)
	require.NoError(t, err)
}
