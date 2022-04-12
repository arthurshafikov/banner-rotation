package handler

import (
	"errors"
	"net/http"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerSlotSocialGroupRoutes(r *router.Router) {
	bannerSlotSocialGroup := r.Group("/")
	{
		bannerSlotSocialGroup.POST("increment", h.incrementClick)
		bannerSlotSocialGroup.GET("getBanner", h.getBannerIDToShow)
	}
}

func (h *Handler) incrementClick(ctx *fasthttp.RequestCtx) {
	bannerID, err := h.getInt64ParamFromRequest(ctx, "banner_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	slotID, err := h.getInt64ParamFromRequest(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupID, err := h.getInt64ParamFromRequest(ctx, "social_group_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlotSocialGroups.IncrementClick(h.ctx, core.IncrementClickInput{
		BannerID:      bannerID,
		SlotID:        slotID,
		SocialGroupID: socialGroupID,
	}); err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}

		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) getBannerIDToShow(ctx *fasthttp.RequestCtx) {
	slotID, err := h.getInt64ParamFromRequest(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupID, err := h.getInt64ParamFromRequest(ctx, "social_group_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	bannerID, err := h.services.BannerSlotSocialGroups.GetBannerIDToShow(h.ctx, core.GetBannerRequest{
		SlotID:        slotID,
		SocialGroupID: socialGroupID,
	})
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}

		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, core.GetBannerResponse{ID: bannerID})
	ctx.SetStatusCode(http.StatusOK)
}
