package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
)

func TestAddBannerSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)

	mock.ExpectQuery("INSERT INTO banner_slots \\(banner_id, slot_id\\) VALUES \\(\\$1, \\$2\\) RETURNING id;").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))

	bannerSlotId, err := bannerSlotRepo.AddBannerSlot(ctx, 1, 2)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(2), bannerSlotId)
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

func TestGetByBannerAndSlotIds(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)
	expected := core.BannerSlot{
		ID:       25,
		BannerId: 1,
		SlotId:   2,
	}

	mock.ExpectQuery("SELECT \\* FROM banner_slots WHERE banner_id=\\$1 AND slot_id=\\$2").
		WithArgs(expected.BannerId, expected.SlotId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "banner_id", "slot_id"}).FromCSVString("25,1,2"))

	result, err := bannerSlotRepo.GetByBannerAndSlotIds(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRandomBannerIdExceptExcluded(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerSlotRepo := NewBannerSlots(sqlxDB)
	expected := core.BannerSlot{
		ID:       25,
		BannerId: 1,
		SlotId:   2,
	}
	excludedBannerSlot := core.BannerSlot{
		ID:       26,
		BannerId: 3,
		SlotId:   2,
	}

	mock.ExpectQuery(`SELECT banner_id FROM banner_slots WHERE slot_id = \$1 
		AND banner_id != \$2 ORDER BY RANDOM\(\) LIMIT 1;`,
	).WithArgs(expected.SlotId, excludedBannerSlot.BannerId).
		WillReturnRows(sqlmock.NewRows([]string{"banner_id"}).FromCSVString("1"))

	result, err := bannerSlotRepo.GetRandomBannerIdExceptExcluded(ctx, expected.SlotId, excludedBannerSlot.BannerId)
	require.NoError(t, err)
	require.Equal(t, expected.BannerId, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
