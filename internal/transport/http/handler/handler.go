package handler

import (
	"context"
	"encoding/json"

	"github.com/arthurshafikov/banner-rotation/internal/services"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type RequestParser interface {
	ParseInt64FromInterface(interface{}) (int64, error)
	ParseInt64FromBytes([]byte) (int64, error)
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

func (h *Handler) setJSONResponse(ctx *fasthttp.RequestCtx, body interface{}) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	ctx.SetBody(bodyJSON)
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
}

func (h *Handler) getIDFromRequest(ctx *fasthttp.RequestCtx) (int64, error) {
	return h.getInt64UserValueFromRequest(ctx, "id")
}

func (h *Handler) getInt64UserValueFromRequest(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	return h.requestParser.ParseInt64FromInterface(ctx.UserValue(key))
}

func (h *Handler) getInt64ParamFromRequest(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	return h.requestParser.ParseInt64FromBytes(ctx.QueryArgs().Peek(key))
}
