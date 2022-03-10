package server

import (
	"github.com/valyala/fasthttp"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve() {
	h := NewHandler()
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			h.home(ctx)
		case "/banner/add":
			h.addBanner(ctx)
		case "/banner/delete":
			h.deleteBanner(ctx)
		case "/banner/get":
			h.getBanner(ctx)
		case "/increment":
			h.incrementClick(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	fasthttp.ListenAndServe(":8123", m)
}
