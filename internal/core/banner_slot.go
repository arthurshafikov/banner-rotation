package core

var BannerSlotTable = "banner_slots"

type BannerSlot struct {
	ID       int64 `db:"id"`
	BannerId int64 `db:"banner_id"`
	SlotId   int64 `db:"slot_id"`
}
