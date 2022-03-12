package http

import (
	"context"
	"fmt"
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

func (s *Server) Serve(port string) {
	r := router.New()
	s.handler.Init(r)

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), r.Handler))
}
