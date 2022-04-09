package handler

import (
	"errors"
	"net/http"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerRoutes(r *router.Router) {
	banners := r.Group("/banner")
	{
		banners.POST("/add", h.addBanner)
		bannersId := banners.Group("/{id:[0-9]+}")
		{
			bannersId.GET("", h.getBanner)
			bannersId.DELETE("/delete", h.deleteBanner)
		}
	}
}

func (h *Handler) addBanner(ctx *fasthttp.RequestCtx) {
	banner := core.Banner{}
	var err error
	banner.ID, err = h.services.Banners.AddBanner(h.ctx, string(ctx.QueryArgs().Peek("description")))
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, banner)
	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) deleteBanner(ctx *fasthttp.RequestCtx) {
	id, err := h.getIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	if err := h.services.Banners.DeleteBanner(h.ctx, id); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) getBanner(ctx *fasthttp.RequestCtx) {
	id, err := h.getIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	banner, err := h.services.Banners.GetBanner(h.ctx, id)
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	h.setJSONResponse(ctx, banner)
}
