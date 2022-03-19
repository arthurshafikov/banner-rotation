package handler

import (
	"errors"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerSlotSocialGroupRoutes(r *router.Router) {
	bannerSlotSocialGroup := r.Group("/")
	{
		bannerSlotSocialGroup.POST("increment", h.incrementClick)
		bannerSlotSocialGroup.GET("getBanner", h.getBannerIdToShow)
	}
}

func (h *Handler) incrementClick(ctx *fasthttp.RequestCtx) {
	bannerId, err := h.getInt64ParamFromRequest(ctx, "banner_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	slotId, err := h.getInt64ParamFromRequest(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupId, err := h.getInt64ParamFromRequest(ctx, "social_group_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlotSocialGroups.IncrementClick(h.ctx, core.IncrementClickInput{
		BannerId:      bannerId,
		SlotId:        slotId,
		SocialGroupId: socialGroupId,
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

func (h *Handler) getBannerIdToShow(ctx *fasthttp.RequestCtx) {
	slotId, err := h.getInt64ParamFromRequest(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupId, err := h.getInt64ParamFromRequest(ctx, "social_group_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	bannerId, err := h.services.BannerSlotSocialGroups.GetBannerIdToShow(h.ctx, core.GetBannerRequest{
		SlotId:        slotId,
		SocialGroupId: socialGroupId,
	})
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}

		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, core.GetBannerResponse{ID: bannerId})
	ctx.SetStatusCode(http.StatusOK)
}
