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
}

func (h *Handler) setJSONResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
}

func (h *Handler) parseIdFromRequest(ctx *fasthttp.RequestCtx) int64 {
	idInterface := ctx.UserValue("id")
	idString, ok := idInterface.(string)
	if !ok {
		panic("not ok wtf")
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	return int64(id)
}
