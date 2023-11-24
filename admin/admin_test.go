package admin

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maxgarvey/cq_server/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error mocking database: %v", err)
	}
	defer mockDB.Close()

	postgres := &postgres.Postgres{
		Connection: mockDB,
	}

	admin := &Admin{
		Postgres: postgres,
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

	mock.ExpectExec(
		"UPDATE cq_server_users " +
			"SET last_login=NOW\\(\\) " +
			"WHERE username=\\$1",
	).WithArgs(expectedUsername).WillReturnResult(sqlmock.NewResult(1, 1))

	loggedIn, err := admin.Login(expectedUsername, expectedPassword)

	assert.Nil(t, err)
	assert.True(t, loggedIn)
}
