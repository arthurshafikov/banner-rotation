package handler

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type RequestParser interface {
	ParseIdFromRequest(*fasthttp.RequestCtx) (int64, error)
	ParseInt64FromRequest(*fasthttp.RequestCtx, string) (int64, error)
	ParseInt64FromQueryArgs(*fasthttp.RequestCtx, string) (int64, error)
}

type Handler struct {
	ctx           context.Context
	services      *services.Services
	requestParser RequestParser
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
	requestParser RequestParser,
) *Handler {
	return &Handler{
		ctx:           ctx,
		services:      services,
		requestParser: requestParser,
	}
}

func (h *Handler) Init(r *router.Router) {
	h.initBannerRoutes(r)
	h.initSlotRoutes(r)
	h.initBannerSlotRoutes(r)
	h.initSocialGroupRoutes(r)
	h.initBannerSlotSocialGroupRoutes(r)
}

func (h *Handler) setJSONResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
}
