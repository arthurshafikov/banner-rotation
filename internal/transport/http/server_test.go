package http

import (
	"context"
	"net/http"
	"testing"

	mock_http "github.com/arthurshafikov/banner-rotation/internal/transport/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestServe(t *testing.T) {
	ctrl := gomock.NewController(t)
	handlerMock := mock_http.NewMockHandler(ctrl)
	server := NewServer(handlerMock)
	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	gomock.InOrder(
		handlerMock.EXPECT().Init(gomock.Any()),
	)
	group.Go(func() error {
		defer cancel()

		response, err := http.Get("http://localhost:9999") //nolint:noctx
		require.NoError(t, response.Body.Close())
		require.NoError(t, err)

		return nil
	})

	server.Serve(ctx, group, "9999")

	require.NoError(t, group.Wait())
}
