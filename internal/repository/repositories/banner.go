package repositories

import "github.com/jmoiron/sqlx"

type Banners struct {
	db *sqlx.DB
}

func NewBanners(db *sqlx.DB) *Banners {
	return &Banners{
		db: db,
	}
}
