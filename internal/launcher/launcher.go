package launcher

import (
	"context"

	"github.com/thewolf27/banner-rotation/internal/repository"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http"
	"github.com/thewolf27/banner-rotation/pkg/postgres"
)

func Launch() {
	ctx := context.TODO()

	db := postgres.NewSqlxDb(ctx, "host=localhost user=homestead password=secret dbname=homestead sslmode=disable")

	repos := repository.NewRepository(db)

	services := services.NewServices(services.Dependencies{
		Repository: repos,
	})

	s := http.NewServer(services)

	s.Serve()
}
