package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserExists(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error mocking database: %v", err)
	}
	defer mockDB.Close()

	postgres := &Postgres{
		Connection: mockDB,
	}

	expectedUsername := "my_user"
	expectedPassword := "my_password"

	expectedCount := 1
	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	mock.ExpectQuery(
		"SELECT COUNT\\(\\*\\) "+
			"FROM cq_server_users "+
			"WHERE username=\\$1 "+
			"AND password=\\$2",
	).WithArgs(expectedUsername, expectedPassword).WillReturnRows(rows)

	exists, err := postgres.UserExists(
		expectedUsername,
		expectedPassword,
	)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestUserDoesNotExist(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error mocking database: %v", err)
	}
	defer mockDB.Close()

	postgres := &Postgres{
		Connection: mockDB,
	}

	expectedUsername := "my_user"
	expectedPassword := "my_password"

	expectedCount := 0
	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	mock.ExpectQuery(
		"SELECT COUNT\\(\\*\\) "+
			"FROM cq_server_users "+
			"WHERE username=\\$1 "+
			"AND password=\\$2",
	).WithArgs(expectedUsername, expectedPassword).WillReturnRows(rows)

	exists, err := postgres.UserExists(
		expectedUsername,
		expectedPassword,
	)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestUpdateLastLogin(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error mocking database: %v", err)
	}
	defer mockDB.Close()

	postgres := &Postgres{
		Connection: mockDB,
	}

	expectedUsername := "my_user"

	mock.ExpectExec(
		"UPDATE cq_server_users " +
			"SET last_login=NOW\\(\\) " +
			"WHERE username=\\$1",
	).WithArgs(expectedUsername).WillReturnResult(sqlmock.NewResult(1, 1))

	err = postgres.UpdateLastLogin(
		expectedUsername,
	)
	assert.Nil(t, err)
}
