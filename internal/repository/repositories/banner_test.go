package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
)

func newSQLXMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB, context.Context) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	return mockDB, mock, sqlx.NewDb(mockDB, "sqlmock"), context.Background()
}

func TestAddBanner(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerRepo := NewBanners(sqlxDB)

	mock.ExpectQuery("INSERT INTO banners \\(description\\) VALUES \\(\\$1\\) RETURNING id;").
		WithArgs("test_description").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))

	bannerId, err := bannerRepo.AddBanner(ctx, "test_description")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(2), bannerId)
}

func TestDeleteBanner(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerRepo := NewBanners(sqlxDB)

	mock.ExpectQuery("DELETE FROM banners WHERE id=\\$1").
		WithArgs(25).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := bannerRepo.DeleteBanner(ctx, 25)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBanner(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	bannerRepo := NewBanners(sqlxDB)
	expected := core.Banner{
		ID:          25,
		Description: "Some banner",
	}

	mock.ExpectQuery("SELECT \\* FROM banners WHERE id=\\$1").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).FromCSVString("25,Some banner"))

	result, err := bannerRepo.GetBanner(ctx, 25)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
