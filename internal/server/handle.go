package server

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) home(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "OK")
}

func (h *Handler) addBanner(ctx *fasthttp.RequestCtx) {
}

func (h *Handler) deleteBanner(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "deleteBanner")
}

func (h *Handler) incrementClick(ctx *fasthttp.RequestCtx) {
}

func (h *Handler) getBanner(ctx *fasthttp.RequestCtx) {
}
