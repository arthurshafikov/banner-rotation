package handler

import (
	"context"
	"encoding/json"
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
		BannerID:      4,
		SlotID:        2,
		SocialGroupID: 1,
	}

	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.BannerID, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SlotID, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SocialGroupID, nil),
		bannerSlotSocialGroupServiceMock.EXPECT().IncrementClick(ctx, input),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.incrementClick, req, resp))
	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetBannerIDToShow(t *testing.T) {
	handler, bannerSlotSocialGroupServiceMock, requestParser, ctx := getMockBannerSlotSocialGroupService(t)
	input := core.GetBannerRequest{
		SlotID:        2,
		SocialGroupID: 1,
	}
	banner := core.Banner{
		ID: 21,
	}

	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SlotID, nil),
		requestParser.EXPECT().ParseInt64FromBytes(nil).Return(input.SocialGroupID, nil),
		bannerSlotSocialGroupServiceMock.EXPECT().GetBannerIDToShow(ctx, input).Return(banner.ID, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.getBannerIDToShow, req, resp))

	expectedJSON, err := json.Marshal(core.GetBannerResponse{ID: banner.ID})
	require.NoError(t, err)
	require.Equal(t, expectedJSON, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
