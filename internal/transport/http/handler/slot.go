package handler

import (
	"errors"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initSlotRoutes(r *router.Router) {
	slots := r.Group("/slot")
	{
		slots.POST("/add", h.addSlot)
		slotsId := slots.Group("/{id:[0-9]+}")
		{
			slotsId.GET("", h.getSlot)
			slotsId.DELETE("/delete", h.deleteSlot)
		}
	}
}

func (h *Handler) addSlot(ctx *fasthttp.RequestCtx) {
	slot := core.Slot{}
	var err error
	slot.ID, err = h.services.Slots.AddSlot(h.ctx, string(ctx.QueryArgs().Peek("description")))
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, slot)
	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) deleteSlot(ctx *fasthttp.RequestCtx) {
	id, err := h.getIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	if err := h.services.Slots.DeleteSlot(h.ctx, id); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) getSlot(ctx *fasthttp.RequestCtx) {
	id, err := h.getIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	slot, err := h.services.Slots.GetSlot(h.ctx, id)
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	h.setJSONResponse(ctx, slot)
}
