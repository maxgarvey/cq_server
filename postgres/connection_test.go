package postgres

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setup() (*sql.DB, sqlmock.Sqlmock, clock.Clock, Postgres) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("Error mocking database: %v", err))
	}

	mockClock := clock.NewMock()
	mockTime, err := time.Parse(
		"Jan 2, 2006 at 3:04pm (MST)",
		"Jan 1, 2020 at 0:00am (PST)",
	)
	mockClock.Set(mockTime)

	postgres := Postgres{
		Connection: mockDB,
		Clock:      mockClock,
	}

	return mockDB, mock, mockClock, postgres
}

func cleanup(mockDB *sql.DB) {
	mockDB.Close()
}

func TestUserExists(t *testing.T) {
	mockDB, mock, _, postgres := setup()
	defer cleanup(mockDB)

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
	mockDB, mock, _, postgres := setup()
	defer cleanup(mockDB)

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
	mockDB, mock, _, postgres := setup()
	defer cleanup(mockDB)

	expectedUsername := "my_user"

	mock.ExpectExec(
		"UPDATE cq_server_users " +
			"SET last_login=NOW\\(\\) " +
			"WHERE username=\\$1",
	).WithArgs(expectedUsername).WillReturnResult(sqlmock.NewResult(1, 1))

	err := postgres.UpdateLastLogin(
		expectedUsername,
	)
	assert.Nil(t, err)
}

func TestCreateSession(t *testing.T) {
	mockDB, mock, mockClock, postgres := setup()
	defer cleanup(mockDB)

	token := fmt.Sprintf("%d%s", 1, mockClock.Now())

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

	session_token, err := postgres.CreateSession(1)
	assert.Nil(t, err)

	decoded_token, err := base64.StdEncoding.DecodeString(session_token)
	if err != nil {
		t.Errorf(
			"Error decoding token: %s",
			err.Error(),
		)
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(decoded_token),
		[]byte(token),
	)
	if err != nil {
		t.Errorf(
			"Hashed token is not a valid bcrypt hash: %v",
			err,
		)
	}
}