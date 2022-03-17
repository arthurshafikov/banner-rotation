package app

import (
	"context"
	"flag"

	"github.com/thewolf27/banner-rotation/internal/config"
	"github.com/thewolf27/banner-rotation/internal/repository"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http"
	"github.com/thewolf27/banner-rotation/internal/transport/http/handler"
	"github.com/thewolf27/banner-rotation/pkg/postgres"
)

var (
	envFileLocation string
)

func init() {
	flag.StringVar(&envFileLocation, "env", "./.env", "Path to .env file")
}

func Run() {
	flag.Parse()

	ctx := context.Background()

	config := config.NewConfig(envFileLocation)

	db := postgres.NewSqlxDb(ctx, config.DSN)

	repos := repository.NewRepository(db)

	services := services.NewServices(services.Dependencies{
		Repository:  repos,
		EGreedValue: config.MultihandedBanditConfig.EGreedValue,
	})

	handler := handler.NewHandler(ctx, services, http.NewRequestParser())
	s := http.NewServer(ctx, handler)
	s.Serve(config.Port)
}
