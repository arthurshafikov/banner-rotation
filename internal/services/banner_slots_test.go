package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_repository "github.com/arthurshafikov/banner-rotation/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getBannerSlotRepoMock(t *testing.T) (context.Context, *mock_repository.MockBannerSlots) {
	t.Helper()
	ctl := gomock.NewController(t)
	bannerSlotRepo := mock_repository.NewMockBannerSlots(ctl)

	return context.Background(), bannerSlotRepo
}

func TestAssociateBannerToSlot(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	expectedBannerSlot := core.BannerSlot{
		ID:       7,
		BannerID: 1,
		SlotID:   2,
	}
	gomock.InOrder(
		bannerSlotRepo.EXPECT().
			AddBannerSlot(ctx, expectedBannerSlot.BannerID, expectedBannerSlot.SlotID).
			Return(expectedBannerSlot.ID, nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	bannerSlotID, err := bs.AssociateBannerToSlot(ctx, expectedBannerSlot.BannerID, expectedBannerSlot.SlotID)

	require.NoError(t, err)
	require.Equal(t, expectedBannerSlot.ID, bannerSlotID)
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

func TestGetByBannerAndSlotIDs(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	expected := &core.BannerSlot{
		ID:       5,
		BannerID: 1,
		SlotID:   2,
	}
	gomock.InOrder(
		bannerSlotRepo.EXPECT().GetByBannerAndSlotIDs(ctx, expected.BannerID, expected.SlotID).
			Return(expected, nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	result, err := bs.GetByBannerAndSlotIDs(ctx, 1, 2)

	require.Equal(t, expected, result)
	require.NoError(t, err)
}

func TestGetRandomBannerIDExceptExcluded(t *testing.T) {
	ctx, bannerSlotRepo := getBannerSlotRepoMock(t)
	expected := &core.BannerSlot{
		ID:       5,
		BannerID: 1,
		SlotID:   3,
	}
	excludedBanner := core.Banner{
		ID: 56,
	}
	gomock.InOrder(
		bannerSlotRepo.EXPECT().
			GetRandomBannerIDExceptExcluded(ctx, expected.SlotID, excludedBanner.ID).
			Return(expected.ID, nil),
	)
	bs := NewBannerSlotService(bannerSlotRepo)

	result, err := bs.GetRandomBannerIDExceptExcluded(ctx, 3, 56)

	require.Equal(t, expected.ID, result)
	require.NoError(t, err)
}
