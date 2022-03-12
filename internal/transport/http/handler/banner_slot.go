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
	bannerId, slotId := h.parseBannerAndSlotIdsFromRequest(ctx)
	if err := h.services.BannerSlots.AssociateBannerToSlot(h.ctx, bannerId, slotId); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) dissociateBannerFromSlot(ctx *fasthttp.RequestCtx) {
	bannerId, slotId := h.parseBannerAndSlotIdsFromRequest(ctx)
	if err := h.services.BannerSlots.DissociateBannerFromSlot(h.ctx, bannerId, slotId); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)

}

func (h *Handler) parseBannerAndSlotIdsFromRequest(ctx *fasthttp.RequestCtx) (int64, int64) {
	return h.parseInt64(ctx.UserValue("bannerId")), h.parseInt64(ctx.UserValue("slotId"))
}
