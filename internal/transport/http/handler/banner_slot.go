package handler

import (
	"net/http"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initBannerSlotRoutes(r *router.Router) {
	bannerSlot := r.Group("/banner/{bannerID:[0-9]+}")
	{
		bannerSlot.POST("/slot/{slotID:[0-9]+}", h.associateBannerToSlot)
		bannerSlot.DELETE("/slot/{slotID:[0-9]+}", h.dissociateBannerFromSlot)
	}
}

func (h *Handler) associateBannerToSlot(ctx *fasthttp.RequestCtx) {
	bannerID, slotID, err := h.parseBannerAndSlotIDsFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	bannerSlot := core.BannerSlot{}
	bannerSlot.ID, err = h.services.BannerSlots.AssociateBannerToSlot(h.ctx, bannerID, slotID)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, bannerSlot)
	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) dissociateBannerFromSlot(ctx *fasthttp.RequestCtx) {
	bannerID, slotID, err := h.parseBannerAndSlotIDsFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	if err := h.services.BannerSlots.DissociateBannerFromSlot(h.ctx, bannerID, slotID); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) parseBannerAndSlotIDsFromRequest(ctx *fasthttp.RequestCtx) (int64, int64, error) {
	bannerID, err := h.getInt64UserValueFromRequest(ctx, "bannerID")
	if err != nil {
		return 0, 0, err
	}

	slotID, err := h.getInt64UserValueFromRequest(ctx, "slotID")
	if err != nil {
		return 0, 0, err
	}

	return bannerID, slotID, nil
}
