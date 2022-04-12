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
	bannerId, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotId, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId, slotId)
	r.NoError(err)
	socialGroupId, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)

	values := url.Values{}
	values.Add("banner_id", fmt.Sprintf("%v", bannerId))
	values.Add("slot_id", fmt.Sprintf("%v", slotId))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupId))

	resp, err := s.makePostRequest("increment", values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE banner_slot_id = $1 AND social_group_id = $2
	`, bannerSlotId, socialGroupId))

	r.Equal(bannerSlotSocialGroup.Clicks, 1)
}

func (s *APITestSuite) TestGetBannerIdToShow() {
	bannerId, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotId, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId, slotId)
	r.NoError(err)
	socialGroupId, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotId))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupId))

	resp, err := s.makeGetRequest("getBanner", values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE banner_slot_id = $1 AND social_group_id = $2
	`, bannerSlotId, socialGroupId))

	r.Equal(bannerSlotSocialGroup.Views, 1)

	expectedResponse := core.GetBannerResponse{ID: bannerId}
	expectedJSON, err := json.Marshal(expectedResponse)
	r.NoError(err)
	bodyJSON, err := ioutil.ReadAll(resp.Body)
	r.NoError(err)
	r.Equal(expectedJSON, bodyJSON)
}

func (s *APITestSuite) TestGetBannerIdToShowWithExistingBannerSlotSocialGroup() {
	bannerId, err := s.repos.Banners.AddBanner(s.ctx, "myBanner")
	r.NoError(err)
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerSlotId, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId, slotId)
	r.NoError(err)
	socialGroupId, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	existingBannerSlotSocialGroup := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId,
		SocialGroupId: socialGroupId,
		Views:         30,
		Clicks:        5,
	}
	err = s.db.QueryRowxContext(
		s.ctx,
		`INSERT INTO banner_slot_social_groups (banner_slot_id, social_group_id, views, clicks) 
			VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		existingBannerSlotSocialGroup.BannerSlotId,
		existingBannerSlotSocialGroup.SocialGroupId,
		existingBannerSlotSocialGroup.Views,
		existingBannerSlotSocialGroup.Clicks,
	).Scan(&existingBannerSlotSocialGroup.ID)
	r.NoError(err)

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotId))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupId))

	resp, err := s.makeGetRequest("getBanner", values)
	r.NoError(err)
	defer resp.Body.Close()
	r.Equal(http.StatusOK, resp.StatusCode)

	bannerSlotSocialGroup := core.BannerSlotSocialGroup{}
	r.NoError(s.db.Get(&bannerSlotSocialGroup, `
		SELECT * FROM banner_slot_social_groups WHERE id = $1
	`, existingBannerSlotSocialGroup.ID))

	r.Equal(bannerSlotSocialGroup.BannerSlotId, existingBannerSlotSocialGroup.BannerSlotId)
	r.Equal(bannerSlotSocialGroup.SocialGroupId, existingBannerSlotSocialGroup.SocialGroupId)
	r.Equal(bannerSlotSocialGroup.Views, existingBannerSlotSocialGroup.Views+1)
	r.Equal(bannerSlotSocialGroup.Clicks, existingBannerSlotSocialGroup.Clicks)
}

func (s *APITestSuite) TestGetBannerIdToShowMultihandBanditIsWorking() {
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerId1, err := s.repos.Banners.AddBanner(s.ctx, "myBanner1")
	r.NoError(err)
	bannerId2, err := s.repos.Banners.AddBanner(s.ctx, "myBanner2")
	r.NoError(err)
	bannerSlotId1, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId1, slotId)
	r.NoError(err)
	bannerSlotId2, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId2, slotId)
	r.NoError(err)
	socialGroupId, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	bannerSlotSocialGroup1 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId1,
		SocialGroupId: socialGroupId,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup1))
	bannerSlotSocialGroup2 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId2,
		SocialGroupId: socialGroupId,
		Views:         5,
		Clicks:        5,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup2))

	bannerShowStatistic := map[int64]int{}

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotId))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupId))

	for i := 0; i < 1000; i++ {
		resp, err := s.makeGetRequest("getBanner", values)
		r.NoError(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		r.NoError(err)
		response := core.GetBannerResponse{}
		r.NoError(json.Unmarshal(body, &response))

		bannerShowStatistic[response.ID]++
	}

	r.True(bannerShowStatistic[bannerId1] < 150 && bannerShowStatistic[bannerId1] > 50)
	r.True(bannerShowStatistic[bannerId2] < 950 && bannerShowStatistic[bannerId2] > 850)
}

func (s *APITestSuite) TestGetBannerIdToShowEveryBannerWasShowed() {
	slotId, err := s.repos.Slots.AddSlot(s.ctx, "mySlot")
	r.NoError(err)
	bannerId1, err := s.repos.Banners.AddBanner(s.ctx, "myBanner1")
	r.NoError(err)
	bannerId2, err := s.repos.Banners.AddBanner(s.ctx, "myBanner2")
	r.NoError(err)
	bannerId3, err := s.repos.Banners.AddBanner(s.ctx, "myBanner3")
	r.NoError(err)
	bannerId4, err := s.repos.Banners.AddBanner(s.ctx, "myBanner4")
	r.NoError(err)
	bannerSlotId1, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId1, slotId)
	r.NoError(err)
	bannerSlotId2, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId2, slotId)
	r.NoError(err)
	bannerSlotId3, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId3, slotId)
	r.NoError(err)
	bannerSlotId4, err := s.repos.BannerSlots.AddBannerSlot(s.ctx, bannerId4, slotId)
	r.NoError(err)
	socialGroupId, err := s.repos.SocialGroups.AddSocialGroup(s.ctx, "mySocialGroup")
	r.NoError(err)
	bannerSlotSocialGroup1 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId1,
		SocialGroupId: socialGroupId,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup1))
	bannerSlotSocialGroup2 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId2,
		SocialGroupId: socialGroupId,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup2))
	bannerSlotSocialGroup3 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId3,
		SocialGroupId: socialGroupId,
		Views:         1000000,
		Clicks:        0,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup3))
	bannerSlotSocialGroup4 := core.BannerSlotSocialGroup{
		BannerSlotId:  bannerSlotId4,
		SocialGroupId: socialGroupId,
		Views:         5,
		Clicks:        5,
	}
	r.NoError(s.insertBannerSlotSocialGroup(&bannerSlotSocialGroup4))

	bannerShowStatistic := map[int64]int{}

	values := url.Values{}
	values.Add("slot_id", fmt.Sprintf("%v", slotId))
	values.Add("social_group_id", fmt.Sprintf("%v", socialGroupId))

	for i := 0; i < 1000; i++ {
		resp, err := s.makeGetRequest("getBanner", values)
		r.NoError(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		r.NoError(err)
		response := core.GetBannerResponse{}
		r.NoError(json.Unmarshal(body, &response))

		bannerShowStatistic[response.ID]++
	}

	r.True(bannerShowStatistic[bannerId1] > 0)
	r.True(bannerShowStatistic[bannerId2] > 0)
	r.True(bannerShowStatistic[bannerId3] > 0)
	r.True(bannerShowStatistic[bannerId4] > bannerShowStatistic[bannerId1]+bannerShowStatistic[bannerId2]+bannerShowStatistic[bannerId3])
}

func (s *APITestSuite) insertBannerSlotSocialGroup(
	bannerSlotSocialGroup *core.BannerSlotSocialGroup,
) error {
	s.T().Helper() //nolint

	return s.db.QueryRowxContext(
		s.ctx,
		`INSERT INTO banner_slot_social_groups (banner_slot_id, social_group_id, views, clicks) 
			VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		bannerSlotSocialGroup.BannerSlotId,
		bannerSlotSocialGroup.SocialGroupId,
		bannerSlotSocialGroup.Views,
		bannerSlotSocialGroup.Clicks,
	).Scan(&bannerSlotSocialGroup.ID)
}
