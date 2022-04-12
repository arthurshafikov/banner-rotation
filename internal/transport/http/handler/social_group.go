package handler

import (
	"errors"
	"net/http"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (h *Handler) initSocialGroupRoutes(r *router.Router) {
	socialGroups := r.Group("/socialGroup")
	{
		socialGroups.POST("/add", h.addSocialGroup)
		socialGroupsID := socialGroups.Group("/{id:[0-9]+}")
		{
			socialGroupsID.GET("", h.getSocialGroup)
			socialGroupsID.DELETE("/delete", h.deleteSocialGroup)
		}
	}
}

func (h *Handler) addSocialGroup(ctx *fasthttp.RequestCtx) {
	socialGroup := core.SocialGroup{}
	var err error
	socialGroup.ID, err = h.services.SocialGroups.AddSocialGroup(h.ctx, string(ctx.QueryArgs().Peek("description")))
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	h.setJSONResponse(ctx, socialGroup)
	ctx.SetStatusCode(http.StatusCreated)
}

func (h *Handler) deleteSocialGroup(ctx *fasthttp.RequestCtx) {
	id, err := h.getIDFromRequest(ctx)
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
	id, err := h.getIDFromRequest(ctx)
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

	ctx.SetStatusCode(http.StatusOK)
	h.setJSONResponse(ctx, socialGroup)
}
