package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/arthurshafikov/banner-rotation/internal/core"
)

func (s *APITestSuite) TestIncrementClick() {
	bannerID, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotID, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID, slotID)
	r.NoError(err)
	socialGroupID, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)

	values := url.Values{}
	values.Add("banner_id", fmt.Sprintf("%v", bannerID))
	values.Add("slot_id", fmt.Sprintf("%v", slotID))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupID))

	resp, err := s.makePostRequest("increment", values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE banner_slot_id = $1 AND social_group_id = $2
	`, bannerSlotID, socialGroupID))

	r.Equal(bannerSlotSocialGroup.Clicks, 1)
}

func (s *APITestSuite) TestGetBannerIDToShow() {
	bannerID, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotID, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID, slotID)
	r.NoError(err)
	socialGroupID, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotID))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupID))

	resp, err := s.makeGetBannerRequest(values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE banner_slot_id = $1 AND social_group_id = $2
	`, bannerSlotID, socialGroupID))

	r.Equal(bannerSlotSocialGroup.Views, 1)

	expectedResponse := core.GetBannerResponse{ID: bannerID}
	expectedJSON, err := json.Marshal(expectedResponse)
	r.NoError(err)
	bodyJSON, err := ioutil.ReadAll(resp.Body)
	r.NoError(err)
	r.Equal(expectedJSON, bodyJSON)
}

func (s *APITestSuite) TestGetBannerIDToShowWithExistingBannerSlotSocialGroup() {
	bannerID, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotID, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID, slotID)
	r.NoError(err)
	socialGroupID, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	existingBannerSlotSocialGroup := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID,
		SocialGroupID: socialGroupID,
		Views:         30,
		Clicks:        5,
	}
	err = s.db.QueryRowxContext(
		s.ctx,
		`INSERT INTO banner_slot_social_groups (banner_slot_id, social_group_id, views, clicks) 
			VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		existingBannerSlotSocialGroup.BannerSlotID,
		existingBannerSlotSocialGroup.SocialGroupID,
		existingBannerSlotSocialGroup.Views,
		existingBannerSlotSocialGroup.Clicks,
	).Scan(&existingBannerSlotSocialGroup.ID)
	r.NoError(err)

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotID))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupID))

	resp, err := s.makeGetBannerRequest(values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE id = $1
	`, existingBannerSlotSocialGroup.ID))

	r.Equal(bannerSlotSocialGroup.BannerSlotID, existingBannerSlotSocialGroup.BannerSlotID)
	r.Equal(bannerSlotSocialGroup.SocialGroupID, existingBannerSlotSocialGroup.SocialGroupID)
	r.Equal(bannerSlotSocialGroup.Views, existingBannerSlotSocialGroup.Views+1)
	r.Equal(bannerSlotSocialGroup.Clicks, existingBannerSlotSocialGroup.Clicks)
}

func (s *APITestSuite) TestGetBannerIDToShowMultihandBanditIsWorking() {
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerID1, err := s.repos.Banners.AddBanner(s.ctx, "myBanner1")
	r.NoError(err)
	bannerID2, err := s.repos.Banners.AddBanner(s.ctx, "myBanner2")
	r.NoError(err)
	bannerSlotID1, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID1, slotID)
	r.NoError(err)
	bannerSlotID2, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID2, slotID)
	r.NoError(err)
	socialGroupID, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	bannerSlotSocialGroup1 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID1,
		SocialGroupID: socialGroupID,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup1))
	bannerSlotSocialGroup2 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID2,
		SocialGroupID: socialGroupID,
		Views:         5,
		Clicks:        5,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup2))

	bannerShowStatistic := map[int64]int{}

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotID))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupID))

	for i := 0; i < 1000; i++ {
		resp, err := s.makeGetBannerRequest(values)
		r.NoError(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		r.NoError(err)
		response := core.GetBannerResponse{}
		r.NoError(json.Unmarshal(body, &response))

		bannerShowStatistic[response.ID]++
	}

	r.True(bannerShowStatistic[bannerID1] < 150 && bannerShowStatistic[bannerID1] > 50)
	r.True(bannerShowStatistic[bannerID2] < 950 && bannerShowStatistic[bannerID2] > 850)
}

func (s *APITestSuite) TestGetBannerIDToShowEveryBannerWasShowed() {
	slotID, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerID1, err := s.repos.Banners.AddBanner(s.ctx, "myBanner1")
	r.NoError(err)
	bannerID2, err := s.repos.Banners.AddBanner(s.ctx, "myBanner2")
	r.NoError(err)
	bannerID3, err := s.repos.Banners.AddBanner(s.ctx, "myBanner3")
	r.NoError(err)
	bannerID4, err := s.repos.Banners.AddBanner(s.ctx, "myBanner4")
	r.NoError(err)
	bannerSlotID1, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID1, slotID)
	r.NoError(err)
	bannerSlotID2, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID2, slotID)
	r.NoError(err)
	bannerSlotID3, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID3, slotID)
	r.NoError(err)
	bannerSlotID4, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerID4, slotID)
	r.NoError(err)
	socialGroupID, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	bannerSlotSocialGroup1 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID1,
		SocialGroupID: socialGroupID,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup1))
	bannerSlotSocialGroup2 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID2,
		SocialGroupID: socialGroupID,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup2))
	bannerSlotSocialGroup3 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID3,
		SocialGroupID: socialGroupID,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup3))
	bannerSlotSocialGroup4 := core.BannerSlotSocialGroup{
		BannerSlotID:  bannerSlotID4,
		SocialGroupID: socialGroupID,
		Views:         5,
		Clicks:        5,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup4))

	bannerShowStatistic := map[int64]int{}

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotID))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupID))

	for i := 0; i < 1000; i++ {
		resp, err := s.makeGetBannerRequest(values)
		r.NoError(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		r.NoError(err)
		response := core.GetBannerResponse{}
		r.NoError(json.Unmarshal(body, &response))

		bannerShowStatistic[response.ID]++
	}

	r.True(bannerShowStatistic[bannerID1] > 0)
	r.True(bannerShowStatistic[bannerID2] > 0)
	r.True(bannerShowStatistic[bannerID3] > 0)
	r.True(
		bannerShowStatistic[bannerID4] >
			bannerShowStatistic[bannerID1]+bannerShowStatistic[bannerID2]+bannerShowStatistic[bannerID3],
	)
}

func (s *APITestSuite) insertBannerSlotSocialGroup(
	bannerSlotSocialGroup *core.BannerSlotSocialGroup,
) error {
	s.T().Helper()

	return s.db.QueryRowxContext(
		s.ctx,
		`INSERT INTO banner_slot_social_groups (banner_slot_id, social_group_id, views, clicks) 
			VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		bannerSlotSocialGroup.BannerSlotID,
		bannerSlotSocialGroup.SocialGroupID,
		bannerSlotSocialGroup.Views,
		bannerSlotSocialGroup.Clicks,
	).Scan(&bannerSlotSocialGroup.ID)
}
