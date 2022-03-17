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

func getMockSocialGroupService(t *testing.T) (*Handler, *mock_services.MockSocialGroups, *mock_handler.MockRequestParser, context.Context) {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	socialGroupServiceMock := mock_services.NewMockSocialGroups(ctrl)
	requestParser := mock_handler.NewMockRequestParser(ctrl)
	handler := NewHandler(ctx, &services.Services{
		SocialGroups: socialGroupServiceMock,
	}, requestParser)

	return handler, socialGroupServiceMock, requestParser, ctx
}

func TestAddSocialGroup(t *testing.T) {
	handler, socialGroupServiceMock, _, ctx := getMockSocialGroupService(t)
	gomock.InOrder(
		socialGroupServiceMock.EXPECT().AddSocialGroup(ctx, "value"),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")
	req.URI().QueryArgs().Add("description", "value")

	resp := fasthttp.AcquireResponse()
	if err := serve(handler.addSocialGroup, req, resp); err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode())
}

func TestDeleteSocialGroup(t *testing.T) {
	handler, socialGroupServiceMock, requestParser, ctx := getMockSocialGroupService(t)
	socialGroupId := int64(20)
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(socialGroupId, nil),
		socialGroupServiceMock.EXPECT().DeleteSocialGroup(ctx, socialGroupId).Return(nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	if err := serve(handler.deleteSocialGroup, req, resp); err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetSocialGroup(t *testing.T) {
	handler, socialGroupServiceMock, requestParser, ctx := getMockSocialGroupService(t)
	socialGroup := core.SocialGroup{
		ID: 20,
	}
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(socialGroup.ID, nil),
		socialGroupServiceMock.EXPECT().GetSocialGroup(ctx, socialGroup.ID).Return(&socialGroup, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	if err := serve(handler.getSocialGroup, req, resp); err != nil {
		t.Error(err)
	}

	expectedSocialGroupJson, err := json.Marshal(socialGroup)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, expectedSocialGroupJson, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
