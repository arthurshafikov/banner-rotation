package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint:gci
)

func NewSqlxDb(ctx context.Context, dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	go func() {
		<-ctx.Done()
		closeConnection(db)
	}()
	if err != nil {
		panic(err)
	}

	return db
}

func closeConnection(db *sqlx.DB) error {
	err := db.Close()

	return err
}
