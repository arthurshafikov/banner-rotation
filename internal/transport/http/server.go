package http

import (
	"context"
	"log"

	"github.com/fasthttp/router"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http/handler"
	"github.com/valyala/fasthttp"
)

type Server struct {
	handler *handler.Handler
}

func NewServer(ctx context.Context, services *services.Services) *Server {
	return &Server{
		handler: handler.NewHandler(ctx, services),
	}
}

func (s *Server) Serve() {
	r := router.New()
	s.handler.Init(r)

	// m := func(ctx *fasthttp.RequestCtx) {
	// 	switch string(ctx.Path()) {
	// 	case "/banner/add":
	// 		s.handler.AddBanner(ctx)
	// 	case "/banner/delete":
	// 		s.handler.DeleteBanner(ctx)
	// 	case "/banner/get":
	// 		s.handler.GetBanner(ctx)
	// 	case "/increment":
	// 		s.handler.IncrementClick(ctx)
	// 	default:
	// 		ctx.Error("not found", fasthttp.StatusNotFound)
	// 	}
	// }

	log.Fatal(fasthttp.ListenAndServe(":8123", r.Handler))
}
