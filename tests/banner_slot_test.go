package tests

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/thewolf27/banner-rotation/internal/core"
)

func (s *APITestSuite) TestAssociateBannerToSlot() {
	bannerId, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)

	resp, err := s.makePostRequest(fmt.Sprintf("/banner/%v/slot/%v", bannerId, slotId), url.Values{})
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusCreated, resp.StatusCode)

	bannerSlot := core.BannerSlot{}
	r.NoError(s.db.Get(&bannerSlot, `
		SELECT * FROM banner_slots WHERE banner_id = $1 AND slot_id = $2
	`, bannerId, slotId))

	r.Equal(bannerSlot.BannerId, bannerId)
	r.Equal(bannerSlot.SlotId, slotId)
}

func (s *APITestSuite) TestDissociateBannerFromSlot() {
	bannerId, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotId, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId, slotId)
	r.NoError(err)

	resp, err := s.makeDeleteRequest(fmt.Sprintf("/banner/%v/slot/%v", bannerId, slotId), url.Values{})
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlot := core.BannerSlot{}
	err = s.db.Get(&bannerSlot, `SELECT * FROM banner_slots WHERE id = $1`, bannerSlotId)
	r.ErrorIs(sql.ErrNoRows, err)
	r.NotEqual(bannerSlotId, bannerSlot.ID)
}
