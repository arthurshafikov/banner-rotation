package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initSlotRoutes(r *router.Router) {
	slots := r.Group("/slot")
	{
		slots.POST("/add", h.AddSlot)
		slotsId := slots.Group("/{id}")
		{
			slotsId.GET("", h.GetSlot)
			slotsId.DELETE("/remove", h.DeleteSlot)
		}
	}
}

func (h *Handler) AddSlot(ctx *fasthttp.RequestCtx) {
	if err := h.services.Slots.AddSlot(h.ctx, string(ctx.QueryArgs().Peek("description"))); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) DeleteSlot(ctx *fasthttp.RequestCtx) {
	id, err := h.parseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	if err := h.services.Slots.DeleteSlot(h.ctx, int64(id)); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) GetSlot(ctx *fasthttp.RequestCtx) {
	id, err := h.parseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	slot, err := h.services.Slots.GetSlot(h.ctx, int64(id))
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}
		ctx.Error(err.Error(), 500)
		return
	}

	slotJSON, err := json.Marshal(slot)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(slotJSON)
	h.setJSONResponse(ctx)
}