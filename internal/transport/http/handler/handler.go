package handler

import (
	"context"
	"errors"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type Handler struct {
	ctx      context.Context
	services *services.Services
}

func NewHandler(ctx context.Context, services *services.Services) *Handler {
	return &Handler{
		ctx:      ctx,
		services: services,
	}
}

func (h *Handler) Init(r *router.Router) {
	h.initBannerRoutes(r)
	h.initSlotRoutes(r)
	h.initBannerSlotRoutes(r)
}

func (h *Handler) setJSONResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
}

func (h *Handler) parseIdFromRequest(ctx *fasthttp.RequestCtx) (int64, error) {
	return h.parseInt64(ctx.UserValue("id"))
}

func (h *Handler) parseInt64(value interface{}) (int64, error) {
	valueString, ok := value.(string)
	if !ok {
		return 0, errors.New("could not convert value to string")
	}

	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return int64(valueInt), nil
}
