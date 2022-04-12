package repositories

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/stretchr/testify/require"
)

func TestAddSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	slotRepo := NewSlots(sqlxDB)

	mock.ExpectQuery("INSERT INTO slots \\(description\\) VALUES \\(\\$1\\) RETURNING id;").
		WithArgs("test_description").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))

	slotId, err := slotRepo.AddSlot(ctx, "test_description")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(2), slotId)
}

func TestDeleteSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	slotRepo := NewSlots(sqlxDB)

	mock.ExpectQuery("DELETE FROM slots WHERE id=\\$1").
		WithArgs(25).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := slotRepo.DeleteSlot(ctx, 25)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSlot(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	slotRepo := NewSlots(sqlxDB)
	expected := core.Slot{
		ID:          25,
		Description: "Some slot",
	}

	mock.ExpectQuery("SELECT \\* FROM slots WHERE id=\\$1").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).FromCSVString("25,Some slot"))

	result, err := slotRepo.GetSlot(ctx, 25)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
