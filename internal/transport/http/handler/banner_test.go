package handler

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/services"
	mock_services "github.com/thewolf27/banner-rotation/internal/services/mocks"
	mock_handler "github.com/thewolf27/banner-rotation/internal/transport/http/handler/mocks"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func getMockBannerService(t *testing.T) (*Handler, *mock_services.MockBanners, *mock_handler.MockRequestParser, context.Context) {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	bannerServiceMock := mock_services.NewMockBanners(ctrl)
	requestParser := mock_handler.NewMockRequestParser(ctrl)
	handler := NewHandler(ctx, &services.Services{
		Banners: bannerServiceMock,
	}, requestParser)

	return handler, bannerServiceMock, requestParser, ctx
}

func TestAddBanner(t *testing.T) {
	handler, bannerServiceMock, _, ctx := getMockBannerService(t)
	expectedBanner := core.Banner{
		ID: 5,
	}
	gomock.InOrder(
		bannerServiceMock.EXPECT().AddBanner(ctx, "value").Return(expectedBanner.ID, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")
	req.URI().QueryArgs().Add("description", "value")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.addBanner, req, resp))
	require.Equal(t, http.StatusCreated, resp.StatusCode())
	expectedJSON, err := json.Marshal(expectedBanner)
	require.NoError(t, err)
	require.Equal(t, expectedJSON, resp.Body())
}

func TestDeleteBanner(t *testing.T) {
	handler, bannerServiceMock, requestParser, ctx := getMockBannerService(t)
	bannerId := int64(20)
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(bannerId, nil),
		bannerServiceMock.EXPECT().DeleteBanner(ctx, bannerId).Return(nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.deleteBanner, req, resp))
	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetBanner(t *testing.T) {
	handler, bannerServiceMock, requestParser, ctx := getMockBannerService(t)
	banner := core.Banner{
		ID: 20,
	}
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(banner.ID, nil),
		bannerServiceMock.EXPECT().GetBanner(ctx, banner.ID).Return(&banner, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.getBanner, req, resp))
	expectedBannerJson, err := json.Marshal(banner)
	require.NoError(t, err)
	require.Equal(t, expectedBannerJson, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func serve(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		if err := fasthttp.Serve(ln, handler); err != nil {
			panic(err)
		}
	}()

	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return client.Do(req, res)
}
