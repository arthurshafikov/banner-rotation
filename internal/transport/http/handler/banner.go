package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerRoutes(r *router.Router) {
	banners := r.Group("/banner")
	{
		banners.POST("/add", h.AddBanner)
		bannersId := banners.Group("/{id}")
		{
			bannersId.GET("", h.GetBanner)
			bannersId.DELETE("/remove", h.DeleteBanner)
		}
	}
}

func (h *Handler) AddBanner(ctx *fasthttp.RequestCtx) {
	if err := h.services.Banners.AddBanner(h.ctx, string(ctx.QueryArgs().Peek("description"))); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) DeleteBanner(ctx *fasthttp.RequestCtx) {
	id, err := h.parseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	if err := h.services.Banners.DeleteBanner(h.ctx, int64(id)); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) GetBanner(ctx *fasthttp.RequestCtx) {
	id, err := h.parseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	banner, err := h.services.Banners.GetBanner(h.ctx, int64(id))
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}
		ctx.Error(err.Error(), 500)
		return
	}

	bannerJSON, err := json.Marshal(banner)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bannerJSON)
	h.setJSONResponse(ctx)
}
