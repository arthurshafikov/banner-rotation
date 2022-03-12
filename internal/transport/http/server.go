package http

import (
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http/handler"
	"github.com/valyala/fasthttp"
)

type Server struct {
	handler *handler.Handler
}

func NewServer(services *services.Services) *Server {
	return &Server{
		handler: handler.NewHandler(services),
	}
}

func (s *Server) Serve() {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/banner/add":
			s.handler.AddBanner(ctx)
		case "/banner/delete":
			s.handler.DeleteBanner(ctx)
		case "/banner/get":
			s.handler.GetBanner(ctx)
		case "/increment":
			s.handler.IncrementClick(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	fasthttp.ListenAndServe(":8123", m)
}
