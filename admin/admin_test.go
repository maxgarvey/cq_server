package admin

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error mocking database: %v", err)
	}
	defer mockDB.Close()
	mockClock := clock.NewMock()
	mockTime, err := time.Parse(
		"Jan 2, 2006 at 3:04pm (MST)",
		"Jan 1, 2020 at 0:00am (PST)",
	)
	mockClock.Set(mockTime)

	postgres := &postgres.Postgres{
		Clock:      mockClock,
		Connection: mockDB,
	}

	admin := &Admin{
		Postgres: postgres,
	}

	expectedUser := data.User{
		ID:        1,
		Username:  "my_user",
		CreatedAt: "",
		LastLogin: "",
	}
	expectedPassword := "my_password"

	// Set up first DB query looking up the user.
	rows := sqlmock.NewRows(
		[]string{"user_id", "username", "created_at", "last_login"},
	).AddRow(
		expectedUser.ID, expectedUser.Username, "", "",
	)
	mock.ExpectQuery(
		"SELECT user_id, username, created_at, last_login "+
			"FROM cq_server_users "+
			"WHERE username=\\$1 "+
			"AND password=\\$2",
	).WithArgs(
		expectedUser.Username, expectedPassword,
	).WillReturnRows(rows)

	// Set up second DB query inserting the new session token.
	mock.ExpectExec("INSERT INTO sessions "+
		"\\(user_id, token, created_at, good_until\\) "+
		"VALUES "+
		"\\(\\$1, \\$2, \\$3, \\$4\\)").WithArgs(
		1,
		sqlmock.AnyArg(),
		mockClock.Now().String(),
		mockClock.Now().Add(time.Hour*24).String(),
	).WillReturnResult(
		sqlmock.NewResult(
			1,
			1,
		),
	)

	// Set up the third DB query, updating the last login for the user.
	mock.ExpectExec(
		"UPDATE cq_server_users " +
			"SET last_login=NOW\\(\\) " +
			"WHERE username=\\$1",
	).WithArgs(
		expectedUser.Username,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := admin.Login(
		expectedUser.Username, expectedPassword,
	)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}
