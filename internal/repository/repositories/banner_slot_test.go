package repositories

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/stretchr/testify/require"
)

func TestAddBannerSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)

	mock.ExpectQuery("INSERT INTO banner_slots \\(banner_id, slot_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id;").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))

	bannerSlotID, err := bannerSlotRepo.AddBannerSlot(ctx, 1, 2)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(2), bannerSlotID)
}

func TestDeleteBannerSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)

	mock.ExpectQuery("DELETE FROM banner_slots WHERE banner_id=\\$1 AND slot_id=\\$2").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerSlotRepo.DeleteBannerSlot(ctx, 1, 2)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByBannerAndSlotIDs(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)
	expected := core.BannerSlot{
		ID:       25,
		BannerID: 1,
		SlotID:   2,
	}

	mock.ExpectQuery("SELECT \\* FROM banner_slots WHERE banner_id=\\$1 AND slot_id=\\$2").
		WithArgs(expected.BannerID, expected.SlotID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "banner_id", "slot_id"}).FromCSVString("25,1,2"))

	result, err := bannerSlotRepo.GetByBannerAndSlotIDs(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRandomBannerIDExceptExcluded(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)
	expected := core.BannerSlot{
		ID:       25,
		BannerID: 1,
		SlotID:   2,
	}
	excludedBannerSlot := core.BannerSlot{
		ID:       26,
		BannerID: 3,
		SlotID:   2,
	}

	mock.ExpectQuery(`SELECT banner_id FROM banner_slots WHERE slot_id = \$1 
		AND banner_id != \$2 ORDER BY RANDOM\(\) LIMIT 1;`,
	).WithArgs(expected.SlotID, excludedBannerSlot.BannerID).
		WillReturnRows(sqlmock.NewRows([]string{"banner_id"}).FromCSVString("1"))

	result, err := bannerSlotRepo.GetRandomBannerIDExceptExcluded(ctx, expected.SlotID, excludedBannerSlot.BannerID)
	require.NoError(t, err)
	require.Equal(t, expected.BannerID, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
