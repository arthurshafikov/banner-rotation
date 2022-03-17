package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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
	slotId, err := h.requestParser.ParseInt64FromQueryArgs(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	bannerId, err := h.requestParser.ParseInt64FromQueryArgs(ctx, "banner_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupId, err := h.requestParser.ParseInt64FromQueryArgs(ctx, "social_group_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlotSocialGroups.IncrementClick(h.ctx, core.IncrementClickInput{
		SlotId:        slotId,
		BannerId:      bannerId,
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
	slotId, err := h.requestParser.ParseInt64FromQueryArgs(ctx, "slot_id")
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroupId, err := h.requestParser.ParseInt64FromQueryArgs(ctx, "social_group_id")
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

	bannerIdJSON, err := json.Marshal(struct{ ID string }{ID: fmt.Sprintf("%v", bannerId)})
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx)
	ctx.SetBody(bannerIdJSON)
	ctx.SetStatusCode(http.StatusOK)
}
