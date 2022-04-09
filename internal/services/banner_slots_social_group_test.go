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
		BannerId: 3,
		SlotId:   1,
	}
	input := core.IncrementClickInput{
		BannerId:      bannerSlot.BannerId,
		SlotId:        bannerSlot.SlotId,
		SocialGroupId: 6,
	}

	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()
	gomock.InOrder(
		bannerSlotService.EXPECT().GetByBannerAndSlotIds(ctx, input.BannerId, input.SlotId).Return(&bannerSlot, nil),
		bannerSlotSocialGroupRepo.EXPECT().IncrementClick(ctx, bannerSlot.ID, input.SocialGroupId).Return(nil),
		queueMock.EXPECT().AddToQueue("clicks", core.IncrementEvent{
			BannerId:      bannerSlot.BannerId,
			SlotId:        bannerSlot.SlotId,
			SocialGroupId: input.SocialGroupId,
			Datetime:      timeNow,
		}).Return(nil),
	)
	bssg := NewBannerSlotSocialGroupService(bannerSlotSocialGroupRepo, bannerSlotService, .1, queueMock)

	err := bssg.IncrementClick(ctx, input)

	require.NoError(t, err)
}

func TestGetBannerIdToShow(t *testing.T) {
	ctx, bannerSlotSocialGroupRepo := getBannerSlotSocialGroupRepoMock(t)
	bannerSlotService := getBannerSlotServiceMock(t)
	queueMock := getQueueMock(t)
	slotId := int64(1)
	mostProfitableBannerSlot := core.BannerSlot{
		ID:       23,
		BannerId: 3,
		SlotId:   slotId,
	}
	randomBannerSlot := core.BannerSlot{
		ID:       44,
		BannerId: 5,
		SlotId:   slotId,
	}
	input := core.GetBannerRequest{
		SlotId:        slotId,
		SocialGroupId: 6,
	}

	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()
	gomock.InOrder(
		bannerSlotSocialGroupRepo.EXPECT().GetTheMostProfitableBannerId(ctx, input.SlotId, input.SocialGroupId).
			Return(mostProfitableBannerSlot.BannerId, nil),
		bannerSlotService.EXPECT().GetRandomBannerIdExceptExcluded(ctx, input.SlotId, mostProfitableBannerSlot.BannerId).
			Return(randomBannerSlot.BannerId, nil),
		bannerSlotService.EXPECT().GetByBannerAndSlotIds(ctx, randomBannerSlot.BannerId, randomBannerSlot.SlotId).
			Return(&randomBannerSlot, nil),
		bannerSlotSocialGroupRepo.EXPECT().IncrementView(ctx, randomBannerSlot.ID, input.SocialGroupId).
			Return(nil),
		queueMock.EXPECT().AddToQueue("views", core.IncrementEvent{
			BannerId:      randomBannerSlot.BannerId,
			SlotId:        randomBannerSlot.SlotId,
			SocialGroupId: input.SocialGroupId,
			Datetime:      timeNow,
		}).Return(nil),
	)
	bssg := NewBannerSlotSocialGroupService(bannerSlotSocialGroupRepo, bannerSlotService, 1, queueMock)

	result, err := bssg.GetBannerIdToShow(ctx, input)

	require.Equal(t, randomBannerSlot.BannerId, result)
	require.NoError(t, err)
}
