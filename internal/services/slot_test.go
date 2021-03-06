package services

import (
	"context"
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_repository "github.com/arthurshafikov/banner-rotation/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getSlotRepoMock(t *testing.T) (context.Context, *mock_repository.MockSlots) {
	t.Helper()
	ctl := gomock.NewController(t)
	slotRepo := mock_repository.NewMockSlots(ctl)

	return context.Background(), slotRepo
}

func TestAddSlot(t *testing.T) {
	ctx, slotRepo := getSlotRepoMock(t)
	expectedSlot := core.Slot{
		ID: 8,
	}
	gomock.InOrder(
		slotRepo.EXPECT().AddSlot(ctx, "test_description").Return(expectedSlot.ID, nil),
	)
	s := NewSlotService(slotRepo)

	slotID, err := s.AddSlot(ctx, "test_description")

	require.NoError(t, err)
	require.Equal(t, expectedSlot.ID, slotID)
}

func TestDeleteSlot(t *testing.T) {
	ctx, slotRepo := getSlotRepoMock(t)
	gomock.InOrder(
		slotRepo.EXPECT().DeleteSlot(ctx, int64(23)).Return(nil),
	)
	s := NewSlotService(slotRepo)

	err := s.DeleteSlot(ctx, 23)

	require.NoError(t, err)
}

func TestGetSlot(t *testing.T) {
	ctx, slotRepo := getSlotRepoMock(t)
	expected := &core.Slot{
		ID: 22,
	}
	gomock.InOrder(
		slotRepo.EXPECT().GetSlot(ctx, expected.ID).Return(expected, nil),
	)
	s := NewSlotService(slotRepo)

	result, err := s.GetSlot(ctx, 22)

	require.Equal(t, expected, result)
	require.NoError(t, err)
}
