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
		BannerSlotID:  2,
		SocialGroupID: 8,
	}

	mock.ExpectQuery("SELECT id FROM banner_slot_social_groups "+
		"WHERE banner_slot_id=\\$1 AND social_group_id=\\$2 LIMIT 1").
		WithArgs(expected.BannerSlotID, expected.SocialGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(fmt.Sprintf("%v", expected.ID)))

	mock.ExpectQuery("UPDATE banner_slot_social_groups SET clicks = clicks \\+ 1 WHERE id = \\$1;").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerSlotSocialGroupRepo.IncrementClick(ctx, expected.BannerSlotID, expected.SocialGroupID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIncrementView(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlotSocialGroup{
		ID:            4,
		BannerSlotID:  2,
		SocialGroupID: 8,
	}

	mock.ExpectQuery("SELECT id FROM banner_slot_social_groups "+
		"WHERE banner_slot_id=\\$1 AND social_group_id=\\$2 LIMIT 1").
		WithArgs(expected.BannerSlotID, expected.SocialGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(fmt.Sprintf("%v", expected.ID)))

	mock.ExpectQuery("UPDATE banner_slot_social_groups SET views = views \\+ 1 WHERE id = \\$1;").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerSlotSocialGroupRepo.IncrementView(ctx, expected.BannerSlotID, expected.SocialGroupID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTheMostProfitableBannerID(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlot{
		ID:       4,
		BannerID: 8,
		SlotID:   9,
	}
	socialGroupID := int64(16)

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
		WithArgs(expected.SlotID, socialGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"banner_id", "ctr"}).FromCSVString(fmt.Sprintf("%v,1", expected.BannerID)))

	bannerID, err := bannerSlotSocialGroupRepo.GetTheMostProfitableBannerID(ctx, expected.SlotID, socialGroupID)
	require.NoError(t, err)
	require.Equal(t, expected.BannerID, bannerID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFirstOrCreateNotExists(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotSocialGroupRepo := NewBannerSlotSocialGroups(sqlxDB)
	expected := core.BannerSlotSocialGroup{
		ID:            4,
		BannerSlotID:  2,
		SocialGroupID: 9,
	}

	mock.ExpectQuery(`SELECT id FROM banner_slot_social_groups 
					  WHERE banner_slot_id=\$1 AND social_group_id=\$2 LIMIT 1`).
		WithArgs(expected.BannerSlotID, expected.SocialGroupID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(`INSERT INTO banner_slot_social_groups 
					  \(banner_slot_id, social_group_id\) VALUES \(\$1, \$2\) RETURNING id;`).
		WithArgs(expected.BannerSlotID, expected.SocialGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("4"))

	result, err := bannerSlotSocialGroupRepo.firstOrCreate(ctx, expected.BannerSlotID, expected.SocialGroupID)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
