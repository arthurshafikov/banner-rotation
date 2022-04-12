package core

import "time"

var BannerSlotSocialGroupTable = "banner_slot_social_groups"

type BannerSlotSocialGroup struct {
	ID            int64 `db:"id"`
	BannerSlotID  int64 `db:"banner_slot_id"`
	SocialGroupID int64 `db:"social_group_id"`
	Views         int   `db:"views"`
	Clicks        int   `db:"clicks"`
}

type IncrementClickInput struct {
	BannerID      int64
	SlotID        int64
	SocialGroupID int64
}

type GetBannerRequest struct {
	SlotID        int64
	SocialGroupID int64
}

type GetBannerResponse struct {
	ID int64 `json:"id"`
}

type IncrementEvent struct {
	BannerID      int64
	SlotID        int64
	SocialGroupID int64
	Datetime      time.Time
}
