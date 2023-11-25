package admin

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	mockClock := clock.NewMock()
	mockTime, err := time.Parse(
		"Jan 2, 2006 at 3:04pm (MST)",
		"Jan 1, 2020 at 0:00am (PST)",
	)
	if err != nil {
		t.Fatalf("error parsing time: %s", err.Error())
	}
	mockClock.Set(mockTime)

	mockPostgres := postgres.InitMock()

	admin := &Admin{
		Clock:    mockClock,
		Postgres: &mockPostgres,
	}

	expectedUser := data.User{
		ID:        1,
		Username:  "my_user",
		CreatedAt: "",
		LastLogin: "",
	}
	expectedPassword := "my_password"

	mockPostgres.SetGetUserResponse(
		postgres.GetUserResponse{
			User:  expectedUser,
			Error: nil,
		},
	)

	token, err := admin.Login(
		expectedUser.Username, expectedPassword,
	)

	// Verify expected response
	assert.Nil(t, err)
	assert.Equal(t, token, "session")

	// Verify mocked DB library calls
	assert.Equal(
		t,
		mockPostgres.GetUserCalls,
		[]postgres.GetUserCall{
			{
				Username: expectedUser.Username,
				Password: expectedPassword,
			},
		},
	)
	assert.Equal(
		t,
		mockPostgres.CreateSessionCalls,
		[]postgres.CreateSessionCall{
			{
				UserID: expectedUser.ID,
			},
		},
	)
	assert.Equal(
		t,
		mockPostgres.UpdateLastLoginCalls,
		[]postgres.UpdateLastLoginCall{
			{
				Username: expectedUser.Username,
			},
		},
	)
}
