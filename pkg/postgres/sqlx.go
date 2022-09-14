package postgres

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint
	"golang.org/x/sync/errgroup"
)

func NewSqlxDB(ctx context.Context, g *errgroup.Group, dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	g.Go(func() error {
		<-ctx.Done()
		return db.Close()
	})

	return db
}
