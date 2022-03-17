package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initSocialGroupRoutes(r *router.Router) {
	socialGroups := r.Group("/socialGroup")
	{
		socialGroups.POST("/add", h.addSocialGroup)
		socialGroupsId := socialGroups.Group("/{id}")
		{
			socialGroupsId.GET("", h.getSocialGroup)
			socialGroupsId.DELETE("/remove", h.deleteSocialGroup)
		}
	}
}

func (h *Handler) addSocialGroup(ctx *fasthttp.RequestCtx) {
	if err := h.services.SocialGroups.AddSocialGroup(h.ctx, string(ctx.QueryArgs().Peek("description"))); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) deleteSocialGroup(ctx *fasthttp.RequestCtx) {
	id, err := h.requestParser.ParseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	if err := h.services.SocialGroups.DeleteSocialGroup(h.ctx, id); err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (h *Handler) getSocialGroup(ctx *fasthttp.RequestCtx) {
	id, err := h.requestParser.ParseIdFromRequest(ctx)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	socialGroup, err := h.services.SocialGroups.GetSocialGroup(h.ctx, id)
	if err != nil {
		if errors.Is(core.ErrNotFound, err) {
			ctx.Error(err.Error(), 404)
			return
		}
		ctx.Error(err.Error(), 500)
		return
	}

	socialGroupJSON, err := json.Marshal(socialGroup)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(socialGroupJSON)
	h.setJSONResponse(ctx)
}
