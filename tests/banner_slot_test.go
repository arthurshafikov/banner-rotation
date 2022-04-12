package tests

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arthurshafikov/banner-rotation/internal/core"
)

func (s *APITestSuite) TestAssociateBannerToSlot() {
	bannerID, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)

	resp, err := s.makePostRequest(fmt.Sprintf("/banner/%v/slot/%v", bannerID, slotID), url.Values{})
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusCreated, resp.StatusCode)

	bannerSlot := core.BannerSlot{}
	r.NoError(s.db.Get(&bannerSlot, `
		SELECT * FROM banner_slots WHERE banner_id = $1 AND slot_id = $2
	`, bannerID, slotID))

	r.Equal(bannerSlot.BannerID, bannerID)
	r.Equal(bannerSlot.SlotID, slotID)
}

func (s *APITestSuite) TestDissociateBannerFromSlot() {
	bannerID, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotID, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID, slotID)
	r.NoError(err)

	resp, err := s.makeDeleteRequest(fmt.Sprintf("/banner/%v/slot/%v", bannerID, slotID), url.Values{})
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlot := core.BannerSlot{}
	err = s.db.Get(&bannerSlot, `SELECT * FROM banner_slots WHERE id = $1`, bannerSlotID)
	r.ErrorIs(sql.ErrNoRows, err)
	r.NotEqual(bannerSlotID, bannerSlot.ID)
}
