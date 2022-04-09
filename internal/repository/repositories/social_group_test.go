package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/stretchr/testify/require"
)

func TestAddSocialGroup(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	socialGroupRepo := NewSocialGroups(sqlxDB)

	mock.ExpectQuery("INSERT INTO social_groups \\(description\\) VALUES \\(\\$1\\) RETURNING id;").
		WithArgs("test_description").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("2"))

	socialGroupId, err := socialGroupRepo.AddSocialGroup(ctx, "test_description")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Equal(t, int64(2), socialGroupId)
}

func TestDeleteSocialGroup(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	socialGroupRepo := NewSocialGroups(sqlxDB)

	mock.ExpectQuery("DELETE FROM social_groups WHERE id=\\$1").
		WithArgs(25).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := socialGroupRepo.DeleteSocialGroup(ctx, 25)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSocialGroup(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	socialGroupRepo := NewSocialGroups(sqlxDB)
	expected := core.SocialGroup{
		ID:          25,
		Description: "Some socialGroup",
	}

	mock.ExpectQuery("SELECT \\* FROM social_groups WHERE id=\\$1").
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).FromCSVString("25,Some socialGroup"))

	result, err := socialGroupRepo.GetSocialGroup(ctx, 25)
	require.NoError(t, err)
	require.Equal(t, &expected, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
