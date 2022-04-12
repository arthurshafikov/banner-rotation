package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint
)

func NewSqlxDB(ctx context.Context, dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	go func() {
		<-ctx.Done()
		if err := closeConnection(db); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}

	return db
}

func closeConnection(db *sqlx.DB) error {
	return db.Close()
}
