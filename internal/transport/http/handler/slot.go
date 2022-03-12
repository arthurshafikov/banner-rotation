package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fasthttp/router"
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
	id := h.parseIdFromRequest(ctx)
	if err := h.services.Slots.DeleteSlot(h.ctx, int64(id)); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) GetSlot(ctx *fasthttp.RequestCtx) {
	id := h.parseIdFromRequest(ctx)
	slot, err := h.services.Slots.GetSlot(h.ctx, int64(id))
	if err != nil {
		// todo sqlx check if empty rows result 404
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
