package repositories

import (
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/stretchr/testify/require"
)

func TestIncrementClick(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlotSocialGroup{
		ID:            4,
		BannerSlotId:  2,
		SocialGroupId: 8,
	}

	mock.ExpectQuery("SELECT id FROM banner_slot_social_groups WHERE banner_slot_id=\\$1 AND social_group_id=\\$2 LIMIT 1").
		WithArgs(expected.BannerSlotId, expected.SocialGroupId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(fmt.Sprintf("%v", expected.ID)))

	mock.ExpectQuery("UPDATE banner_slot_social_groups SET clicks = clicks \\+ 1 WHERE id = \\$1;").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerSlotSocialGroupRepo.IncrementClick(ctx, expected.BannerSlotId, expected.SocialGroupId)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIncrementView(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlotSocialGroup{
		ID:            4,
		BannerSlotId:  2,
		SocialGroupId: 8,
	}

	mock.ExpectQuery("SELECT id FROM banner_slot_social_groups WHERE banner_slot_id=\\$1 AND social_group_id=\\$2 LIMIT 1").
		WithArgs(expected.BannerSlotId, expected.SocialGroupId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(fmt.Sprintf("%v", expected.ID)))

	mock.ExpectQuery("UPDATE banner_slot_social_groups SET views = views \\+ 1 WHERE id = \\$1;").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerSlotSocialGroupRepo.IncrementView(ctx, expected.BannerSlotId, expected.SocialGroupId)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTheMostProfitableBannerId(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlot{
		ID:       4,
		BannerId: 8,
		SlotId:   9,
	}
	socialGroupId := int64(16)

	expectedQuery := fmt.Sprintf(
		`SELECT %[2]s\.banner_id, CAST\(%[1]s\.clicks AS DECIMAL\)\/(.+?) as ctr FROM %[1]s
			LEFT JOIN %[2]s ON %[2]s\.id = %[1]s\.banner_slot_id
			WHERE %[2]s\.slot_id = \$1 AND %[1]s\.social_group_id = \$2
			ORDER BY ctr DESC
			LIMIT 1;`,
		core.BannerSlotSocialGroupTable,
		core.BannerSlotTable,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(expected.SlotId, socialGroupId).
		WillReturnRows(sqlmock.NewRows([]string{"banner_id", "ctr"}).FromCSVString(fmt.Sprintf("%v,1", expected.BannerId)))

	bannerId, err := bannerSlotSocialGroupRepo.GetTheMostProfitableBannerId(ctx, expected.SlotId, socialGroupId)
	require.NoError(t, err)
	require.Equal(t, expected.BannerId, bannerId)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFirstOrCreateNotExists(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlotSocialGroup{
		ID:            4,
		BannerSlotId:  2,
		SocialGroupId: 9,
	}

	mock.ExpectQuery(`SELECT id FROM banner_slot_social_groups 
					  WHERE banner_slot_id=\$1 AND social_group_id=\$2 LIMIT 1`).
		WithArgs(expected.BannerSlotId, expected.SocialGroupId).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(`INSERT INTO banner_slot_social_groups 
					  \(banner_slot_id, social_group_id\) VALUES \(\$1, \$2\) RETURNING id;`).
		WithArgs(expected.BannerSlotId, expected.SocialGroupId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("4"))

	result, err := bannerSlotSocialGroupRepo.firstOrCreate(ctx, expected.BannerSlotId, expected.SocialGroupId)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
