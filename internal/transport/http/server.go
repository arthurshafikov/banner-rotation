package http

import (
	"context"
	"fmt"
	"log"

	"github.com/arthurshafikov/banner-rotation/internal/transport/http/handler"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	handler *handler.Handler
	server  *fasthttp.Server
}

func NewServer(handler *handler.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Serve(ctx context.Context, g *errgroup.Group, port string) {
	r := router.New()
	s.handler.Init(r)

	s.server = &fasthttp.Server{
		Handler: r.Handler,
	}

	g.Go(func() error {
		<-ctx.Done()

		s.shutdown()

		return nil
	})

	g.Go(func() error {
		return s.server.ListenAndServe(fmt.Sprintf(":%s", port))
	})
}

func (s *Server) shutdown() {
	log.Println("Shutdown Server ...")

	if err := s.server.Shutdown(); err != nil {
		log.Println("Server forced to shutdown: ", err)
	}
}
