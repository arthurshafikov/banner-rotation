package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/arthurshafikov/banner-rotation/internal/services"
	mock_services "github.com/arthurshafikov/banner-rotation/internal/services/mocks"
	mock_handler "github.com/arthurshafikov/banner-rotation/internal/transport/http/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func getMockBannerSlotService(t *testing.T) (*Handler, *mock_services.MockBannerSlots, *mock_handler.MockRequestParser, context.Context) {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	bannerSlotServiceMock := mock_services.NewMockBannerSlots(ctrl)
	requestParser := mock_handler.NewMockRequestParser(ctrl)
	handler := NewHandler(ctx, &services.Services{
		BannerSlots: bannerSlotServiceMock,
	}, requestParser)

	return handler, bannerSlotServiceMock, requestParser, ctx
}

func TestAssociateBannerToSlot(t *testing.T) {
	handler, bannerSlotServiceMock, requestParser, ctx := getMockBannerSlotService(t)
	bannerSlot := core.BannerSlot{
		BannerId: 4,
		SlotId:   8,
	}
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(bannerSlot.BannerId, nil),
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(bannerSlot.SlotId, nil),
		bannerSlotServiceMock.EXPECT().AssociateBannerToSlot(ctx, bannerSlot.BannerId, bannerSlot.SlotId),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.associateBannerToSlot, req, resp))
	require.Equal(t, http.StatusCreated, resp.StatusCode())
}

func TestDissociateBannerFromSlot(t *testing.T) {
	handler, bannerSlotServiceMock, requestParser, ctx := getMockBannerSlotService(t)
	bannerSlot := core.BannerSlot{
		BannerId: 4,
		SlotId:   8,
	}
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(bannerSlot.BannerId, nil),
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(bannerSlot.SlotId, nil),
		bannerSlotServiceMock.EXPECT().DissociateBannerFromSlot(ctx, bannerSlot.BannerId, bannerSlot.SlotId),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.dissociateBannerFromSlot, req, resp))
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
