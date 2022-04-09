package http

import (
	"context"
	"fmt"
	"log"

	"github.com/arthurshafikov/banner-rotation/internal/transport/http/handler"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Server struct {
	handler *handler.Handler
}

func NewServer(ctx context.Context, handler *handler.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Serve(port string) {
	r := router.New()
	s.handler.Init(r)

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), r.Handler))
}
