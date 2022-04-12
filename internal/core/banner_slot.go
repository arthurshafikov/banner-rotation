package core

var BannerSlotTable = "banner_slots"

type BannerSlot struct {
	ID       int64 `db:"id"`
	BannerID int64 `db:"banner_id"`
	SlotID   int64 `db:"slot_id"`
}
