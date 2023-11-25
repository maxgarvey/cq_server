package admin

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/postgres"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (clock.Clock, *postgres.MockPostgres, *Admin) {
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

	return mockClock, &mockPostgres, admin
}

func TestValidateSession(t *testing.T) {
	mockClock, mockPostgres, admin := setup(t)

	mockPostgres.GetSessionResponse = postgres.GetSessionResponse{
		Session: data.Session{
			UserID:    1,
			Token:     "session",
			CreatedAt: mockClock.Now(),
			GoodUntil: mockClock.Now().Add(time.Hour * 24),
		},
		Error: nil,
	}

	valid, err := admin.ValidateSession("session")

	// Verify response
	assert.Equal(t, valid, true)
	assert.Nil(t, err)

	// Verify mocked DB call
	assert.Equal(
		t,
		mockPostgres.GetSessionCalls,
		[]postgres.GetSessionCall{
			{
				Token: "session",
			},
		},
	)
}

func TestValidateBadSession(t *testing.T) {
	mockClock, mockPostgres, admin := setup(t)

	mockPostgres.GetSessionResponse = postgres.GetSessionResponse{
		Session: data.Session{
			UserID:    1,
			Token:     "session",
			CreatedAt: mockClock.Now(),
			GoodUntil: mockClock.Now().Add(time.Hour * -24),
		},
		Error: nil,
	}

	valid, err := admin.ValidateSession("session")

	// Verify response
	assert.Equal(t, valid, false)
	assert.Nil(t, err)
}

func TestExtendSession(t *testing.T) {
	_, mockPostgres, admin := setup(t)

	mockPostgres.SetExtendSessionResponse(
		postgres.ExtendSessionResponse{
			Error: nil,
		},
	)

	err := admin.ExtendSession("session")

	assert.Nil(t, err)
}
