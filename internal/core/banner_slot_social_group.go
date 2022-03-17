package core

var BannerSlotSocialGroupTable = "banner_slot_social_groups"

type BannerSlotSocialGroup struct {
	ID            int64 `db:"id"`
	BannerSlotId  int64 `db:"banner_slot_id"`
	SocialGroupId int64 `db:"social_group_id"`
	Views         int   `db:"views"`
	Clicks        int   `db:"clicks"`
}

type IncrementClickInput struct {
	BannerId      int64
	SlotId        int64
	SocialGroupId int64
}

type GetBannerRequest struct {
	SlotId        int64
	SocialGroupId int64
}

type GetBannerResponse struct {
	ID int64 `json:"id"`
}
