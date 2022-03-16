package core

var BannersTable = "banners"

type Banner struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}
