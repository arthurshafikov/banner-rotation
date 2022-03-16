package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/thewolf27/banner-rotation/internal/core"
)

func TestAddSocialGroup(t *testing.T) {
	mockDB, mock, sqlxDB, ctx := newSQLXMock(t)
	defer mockDB.Close()
	socialGroupRepo := NewSocialGroups(sqlxDB)

	mock.ExpectQuery("INSERT INTO social_groups \\(description\\) VALUES \\(\\$1\\)").
		WithArgs("test_description").
		WillReturnRows(sqlmock.NewRows([]string{}))

	err := socialGroupRepo.AddSocialGroup(ctx, "test_description")
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
