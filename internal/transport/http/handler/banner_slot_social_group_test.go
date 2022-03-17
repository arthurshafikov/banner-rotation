package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/services"
	mock_services "github.com/thewolf27/banner-rotation/internal/services/mocks"
	mock_handler "github.com/thewolf27/banner-rotation/internal/transport/http/handler/mocks"
	"github.com/valyala/fasthttp"
)

func getMockBannerSlotSocialGroupService(
	t *testing.T,
) (*Handler, *mock_services.MockBannerSlotSocialGroups, *mock_handler.MockRequestParser, context.Context) {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	bannerSlotSocialGroupServiceMock := mock_services.NewMockBannerSlotSocialGroups(ctrl)
	requestParser := mock_handler.NewMockRequestParser(ctrl)
	handler := NewHandler(ctx, &services.Services{
		BannerSlotSocialGroups: bannerSlotSocialGroupServiceMock,
	}, requestParser)

	return handler, bannerSlotSocialGroupServiceMock, requestParser, ctx
}

func TestIncrementClick(t *testing.T) {
	handler, bannerSlotSocialGroupServiceMock, requestParser, ctx := getMockBannerSlotSocialGroupService(t)
	input := core.IncrementClickInput{
		BannerId:      4,
		SlotId:        2,
		SocialGroupId: 1,
	}

	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.BannerId, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SlotId, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SocialGroupId, nil),
		bannerSlotSocialGroupServiceMock.EXPECT().IncrementClick(ctx, input),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	if err := serve(handler.incrementClick, req, resp); err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetBannerIdToShow(t *testing.T) {
	handler, bannerSlotSocialGroupServiceMock, requestParser, ctx := getMockBannerSlotSocialGroupService(t)
	input := core.GetBannerRequest{
		SlotId:        2,
		SocialGroupId: 1,
	}
	banner := core.Banner{
		ID: 21,
	}

	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SlotId, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SocialGroupId, nil),
		bannerSlotSocialGroupServiceMock.EXPECT().GetBannerIdToShow(ctx, input).Return(banner.ID, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	if err := serve(handler.getBannerIdToShow, req, resp); err != nil {
		t.Error(err)
	}

	expectedJson, err := json.Marshal(core.GetBannerResponse{ID: banner.ID})
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, expectedJson, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
