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

func getMockSlotService(
	t *testing.T,
) (*Handler, *mock_services.MockSlots, *mock_handler.MockRequestParser, context.Context) {
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
	expectedSlot := core.Slot{
		ID: 8,
	}
	gomock.InOrder(
		slotServiceMock.EXPECT().AddSlot(ctx, "value").Return(expectedSlot.ID, nil),
	)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost")
	req.URI().QueryArgs().Add("description", "value")

	resp := fasthttp.AcquireResponse()
	require.NoError(t, serve(handler.addSlot, req, resp))
	require.Equal(t, http.StatusCreated, resp.StatusCode())
	expectedJSON, err := json.Marshal(expectedSlot)
	require.NoError(t, err)
	require.Equal(t, expectedJSON, resp.Body())
}

func TestDeleteSlot(t *testing.T) {
	handler, slotServiceMock, requestParser, ctx := getMockSlotService(t)
	slotID := int64(20)
	gomock.InOrder(
		requestParser.EXPECT().ParseInt64FromInterface(nil).Return(slotID, nil),
		slotServiceMock.EXPECT().DeleteSlot(ctx, slotID).Return(nil),
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

	expectedSlotJSON, err := json.Marshal(slot)
	require.NoError(t, err)
	require.Equal(t, expectedSlotJSON, resp.Body())
	require.Equal(t, http.StatusOK, resp.StatusCode())
}
