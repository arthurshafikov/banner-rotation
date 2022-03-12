package handler

import (
	"context"
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

func (h *Handler) parseIdFromRequest(ctx *fasthttp.RequestCtx) int64 {
	return h.parseInt64(ctx.UserValue("id"))
}

func (h *Handler) parseInt64(num interface{}) int64 {
	numString, ok := num.(string)
	if !ok {
		panic("not ok wtf")
	}

	numInt, err := strconv.Atoi(numString)
	if err != nil {
		panic(err)
	}

	return int64(numInt)
}
