package services

import (
	"context"
	"testing"
	"time"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	mock_repository "github.com/arthurshafikov/banner-rotation/internal/repository/mocks"
	mock_services "github.com/arthurshafikov/banner-rotation/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tkuchiki/faketime"
)

func getBannerSlotSocialGroupRepoMock(t *testing.T) (context.Context, *mock_repository.MockBannerSlotSocialGroups) {
	t.Helper()
	ctl := gomock.NewController(t)
	bannerSlotSocialGroupRepo := mock_repository.NewMockBannerSlotSocialGroups(ctl)

	return context.Background(), bannerSlotSocialGroupRepo
}

func getBannerSlotServiceMock(t *testing.T) *mock_services.MockBannerSlots {
	t.Helper()
	ctl := gomock.NewController(t)

	return mock_services.NewMockBannerSlots(ctl)
}

func getQueueMock(t *testing.T) *mock_services.MockQueue {
	t.Helper()
	ctl := gomock.NewController(t)

	return mock_services.NewMockQueue(ctl)
}

func TestIncrementClick(t *testing.T) {
	ctx, bannerSlotSocialGroupRepo := getBannerSlotSocialGroupRepoMock(t)
	bannerSlotService := getBannerSlotServiceMock(t)
	queueMock := getQueueMock(t)
	bannerSlot := core.BannerSlot{
		ID:       23,
		BannerID: 3,
		SlotID:   1,
	}
	input := core.IncrementClickInput{
		BannerID:      bannerSlot.BannerID,
		SlotID:        bannerSlot.SlotID,
		SocialGroupID: 6,
	}

	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()
	gomock.InOrder(
		bannerSlotService.EXPECT().GetByBannerAndSlotIDs(ctx, input.BannerID, input.SlotID).Return(&bannerSlot, nil),
		bannerSlotSocialGroupRepo.EXPECT().IncrementClick(ctx, bannerSlot.ID, input.SocialGroupID).Return(nil),
		queueMock.EXPECT().AddToQueue("clicks", core.IncrementEvent{
			BannerID:      bannerSlot.BannerID,
			SlotID:        bannerSlot.SlotID,
			SocialGroupID: input.SocialGroupID,
			Datetime:      timeNow,
		}).Return(nil),
	)
	bssg := NewBannerSlotSocialGroupService(bannerSlotSocialGroupRepo, bannerSlotService, .1, queueMock)

	err := bssg.IncrementClick(ctx, input)

	require.NoError(t, err)
}

func TestGetBannerIDToShow(t *testing.T) {
	ctx, bannerSlotSocialGroupRepo := getBannerSlotSocialGroupRepoMock(t)
	bannerSlotService := getBannerSlotServiceMock(t)
	queueMock := getQueueMock(t)
	slotID := int64(1)
	mostProfitableBannerSlot := core.BannerSlot{
		ID:       23,
		BannerID: 3,
		SlotID:   slotID,
	}
	randomBannerSlot := core.BannerSlot{
		ID:       44,
		BannerID: 5,
		SlotID:   slotID,
	}
	input := core.GetBannerRequest{
		SlotID:        slotID,
		SocialGroupID: 6,
	}

	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()
	gomock.InOrder(
		bannerSlotSocialGroupRepo.EXPECT().GetTheMostProfitableBannerID(ctx, input.SlotID, input.SocialGroupID).
			Return(mostProfitableBannerSlot.BannerID, nil),
		bannerSlotService.EXPECT().GetRandomBannerIDExceptExcluded(ctx, input.SlotID, mostProfitableBannerSlot.BannerID).
			Return(randomBannerSlot.BannerID, nil),
		bannerSlotService.EXPECT().GetByBannerAndSlotIDs(ctx, randomBannerSlot.BannerID, randomBannerSlot.SlotID).
			Return(&randomBannerSlot, nil),
		bannerSlotSocialGroupRepo.EXPECT().IncrementView(ctx, randomBannerSlot.ID, input.SocialGroupID).
			Return(nil),
		queueMock.EXPECT().AddToQueue("views", core.IncrementEvent{
			BannerID:      randomBannerSlot.BannerID,
			SlotID:        randomBannerSlot.SlotID,
			SocialGroupID: input.SocialGroupID,
			Datetime:      timeNow,
		}).Return(nil),
	)
	bssg := NewBannerSlotSocialGroupService(bannerSlotSocialGroupRepo, bannerSlotService, 1, queueMock)

	result, err := bssg.GetBannerIDToShow(ctx, input)

	require.Equal(t, randomBannerSlot.BannerID, result)
	require.NoError(t, err)
}
