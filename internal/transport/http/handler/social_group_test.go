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

func getMockSocialGroupService(
	t *testing.T,
) (*Handler, *mock_services.MockSocialGroups, *mock_handler.MockRequestParser, context.Context) {
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
	expectedSocialGroup := core.SocialGroup{
		ID: 5,
	}
	gomock.InOrder(
		socialGroupServiceMock.EXPECT().AddSocialGroup(ctx, "value").Return(expectedSocialGroup.ID, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")
	req.URI().QueryArgs().Add("description", "value")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.addSocialGroup, req, resp))
	require.Equal(t, http.StatusCreated, resp.StatusCode())
	expectedJSON, err := json.Marshal(expectedSocialGroup)
	require.NoError(t, err)
	require.Equal(t, expectedJSON, resp.Body())
}

func TestDeleteSocialGroup(t *testing.T) {
	handler, socialGroupServiceMock, requestParser, ctx := getMockSocialGroupService(t)
	socialGroupID := int64(20)
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(socialGroupID, nil),
		socialGroupServiceMock.EXPECT().DeleteSocialGroup(ctx, socialGroupID).Return(nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.deleteSocialGroup, req, resp))
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
	require.NoError(t, serve(handler.getSocialGroup, req, resp))

	expectedSocialGroupJSON, err := json.Marshal(socialGroup)
	require.NoError(t, err)
	require.Equal(t, expectedSocialGroupJSON, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
