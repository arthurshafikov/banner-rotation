package launcher

import (
	"context"
	"flag"

	"github.com/thewolf27/banner-rotation/internal/config"
	"github.com/thewolf27/banner-rotation/internal/repository"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http"
	"github.com/thewolf27/banner-rotation/pkg/postgres"
)

var (
	envFileLocation string
)

func init() {
	flag.StringVar(&envFileLocation, "env", "./.env", "Path to .env file")
}

func Launch() {
	flag.Parse()

	ctx := context.Background()

	config := config.NewConfig(envFileLocation)

	db := postgres.NewSqlxDb(ctx, config.DSN)

	repos := repository.NewRepository(db)

	services := services.NewServices(services.Dependencies{
		Repository: repos,
	})

	s := http.NewServer(ctx, services)

	s.Serve(config.Port)
}
