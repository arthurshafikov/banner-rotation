package handler

import (
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerSlotRoutes(r *router.Router) {
	bannerSlot := r.Group("/banner/{bannerId}")
	{
		bannerSlot.POST("/slot/{slotId}", h.associateBannerToSlot)
		bannerSlot.DELETE("/slot/{slotId}", h.dissociateBannerFromSlot)
	}
}

func (h *Handler) associateBannerToSlot(ctx *fasthttp.RequestCtx) {
	bannerId, slotId, err := h.parseBannerAndSlotIdsFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlots.AssociateBannerToSlot(h.ctx, bannerId, slotId); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) dissociateBannerFromSlot(ctx *fasthttp.RequestCtx) {
	bannerId, slotId, err := h.parseBannerAndSlotIdsFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlots.DissociateBannerFromSlot(h.ctx, bannerId, slotId); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)

}

func (h *Handler) parseBannerAndSlotIdsFromRequest(ctx *fasthttp.RequestCtx) (int64, int64, error) {
	bannerId, err := h.parseInt64FromInterface(ctx.UserValue("bannerId"))
	if err != nil {
		return 0, 0, err
	}

	slotId, err := h.parseInt64FromInterface(ctx.UserValue("slotId"))
	if err != nil {
		return 0, 0, err
	}

	return bannerId, slotId, nil
}
