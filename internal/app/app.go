package app

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/arthurshafikov/banner-rotation/internal/config"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
	"github.com/arthurshafikov/banner-rotation/internal/services"
	"github.com/arthurshafikov/banner-rotation/internal/transport/http"
	"github.com/arthurshafikov/banner-rotation/internal/transport/http/handler"
	"github.com/arthurshafikov/banner-rotation/pkg/postgres"
	"github.com/arthurshafikov/banner-rotation/pkg/queue"
)

var (
	envFileLocation string
)

func init() {
	flag.StringVar(&envFileLocation, "env", "./.env", "Path to .env file")
}

func Run() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(envFileLocation)

	db := postgres.NewSqlxDb(ctx, config.DSN)

	repos := repository.NewRepository(db)

	queue := queue.NewQueue(ctx, config.QueueConfig.BrokerAddress)
	go queue.Dispatch()

	services := services.NewServices(services.Dependencies{
		Repository:  repos,
		EGreedValue: config.MultihandedBanditConfig.EGreedValue,
		Queue:       queue,
	})

	handler := handler.NewHandler(ctx, services, http.NewRequestParser())
	s := http.NewServer(ctx, handler)
	s.Serve(config.ServerConfig.Port)
}
