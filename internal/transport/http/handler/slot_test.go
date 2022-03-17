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

func getMockSlotService(t *testing.T) (*Handler, *mock_services.MockSlots, *mock_handler.MockRequestParser, context.Context) {
	t.Helper()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	slotServiceMock := mock_services.NewMockSlots(ctrl)
	requestParser := mock_handler.NewMockRequestParser(ctrl)
	handler := NewHandler(ctx, &services.Services{
		Slots: slotServiceMock,
	}, requestParser)

	return handler, slotServiceMock, requestParser, ctx
}

func TestAddSlot(t *testing.T) {
	handler, slotServiceMock, _, ctx := getMockSlotService(t)
	gomock.InOrder(
		slotServiceMock.EXPECT().AddSlot(ctx, "value"),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")
	req.URI().QueryArgs().Add("description", "value")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.addSlot, req, resp))
	require.Equal(t, http.StatusCreated, resp.StatusCode())
}

func TestDeleteSlot(t *testing.T) {
	handler, slotServiceMock, requestParser, ctx := getMockSlotService(t)
	slotId := int64(20)
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(slotId, nil),
		slotServiceMock.EXPECT().DeleteSlot(ctx, slotId).Return(nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.deleteSlot, req, resp))
	require.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestGetSlot(t *testing.T) {
	handler, slotServiceMock, requestParser, ctx := getMockSlotService(t)
	slot := core.Slot{
		ID: 20,
	}
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(slot.ID, nil),
		slotServiceMock.EXPECT().GetSlot(ctx, slot.ID).Return(&slot, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.getSlot, req, resp))

	expectedSlotJson, err := json.Marshal(slot)
	require.NoError(t, err)
	require.Equal(t, expectedSlotJson, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
