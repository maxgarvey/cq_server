package postgres

import (
	"time"

	"github.com/maxgarvey/cq_server/data"
)

type MockPostgres struct {
	GetUserCalls            []GetUserCall
	GetUserResponse         GetUserResponse
	UpdateLastLoginCalls    []UpdateLastLoginCall
	UpdateLastLoginResponse UpdateLastLoginResponse
	CreateSessionCalls      []CreateSessionCall
	CreateSessionResponse   CreateSessionResponse
	GetSessionCalls         []GetSessionCall
	GetSessionResponse      GetSessionResponse
	ExtendSessionCalls      []ExtendSessionCall
	ExtendSessionResponse   ExtendSessionResponse
}

func InitMock() MockPostgres {
	return MockPostgres{
		GetUserCalls: []GetUserCall{},
		GetUserResponse: GetUserResponse{
			User:  data.User{},
			Error: nil,
		},
		UpdateLastLoginCalls: []UpdateLastLoginCall{},
		UpdateLastLoginResponse: UpdateLastLoginResponse{
			Error: nil,
		},
		CreateSessionCalls: []CreateSessionCall{},
		CreateSessionResponse: CreateSessionResponse{
			Token: "session",
			Error: nil,
		},
		GetSessionCalls: []GetSessionCall{},
		GetSessionResponse: GetSessionResponse{
			Session: data.Session{
				UserID:    1,
				Token:     "session",
				CreatedAt: time.Now(),
				GoodUntil: time.Now(),
			},
			Error: nil,
		},
		ExtendSessionCalls: []ExtendSessionCall{},
		ExtendSessionResponse: ExtendSessionResponse{
			Error: nil,
		},
	}
}

type GetUserCall struct {
	Username string
	Password string
}

type GetUserResponse struct {
	User  data.User
	Error error
}

func (m *MockPostgres) SetGetUserResponse(response GetUserResponse) {
	m.GetUserResponse = response
}

func (m *MockPostgres) GetUser(username string, password string) (data.User, error) {
	m.GetUserCalls = append(
		m.GetUserCalls,
		GetUserCall{
			Username: username,
			Password: password,
		},
	)
	return m.GetUserResponse.User, m.GetUserResponse.Error
}

type UpdateLastLoginCall struct {
	Username string
}

type UpdateLastLoginResponse struct {
	Error error
}

func (m *MockPostgres) SetUpdateLastLoginResponse(
	response UpdateLastLoginResponse,
) {
	m.UpdateLastLoginResponse = response
}

func (m *MockPostgres) UpdateLastLogin(username string) error {
	m.UpdateLastLoginCalls = append(
		m.UpdateLastLoginCalls,
		UpdateLastLoginCall{
			Username: username,
		},
	)
	return m.UpdateLastLoginResponse.Error
}

type CreateSessionCall struct {
	UserID int
}

type CreateSessionResponse struct {
	Token string
	Error error
}

func (m *MockPostgres) SetCreateSessionResponse(response CreateSessionResponse) {
	m.CreateSessionResponse = response
}

func (m *MockPostgres) CreateSession(user_id int) (string, error) {
	m.CreateSessionCalls = append(
		m.CreateSessionCalls,
		CreateSessionCall{
			UserID: user_id,
		},
	)
	return m.CreateSessionResponse.Token, m.CreateSessionResponse.Error
}

type GetSessionCall struct {
	Token string
}

type GetSessionResponse struct {
	Session data.Session
	Error   error
}

func (m *MockPostgres) SetGetSessionResponse(response GetSessionResponse) {
	m.GetSessionResponse = response
}

func (m *MockPostgres) GetSession(token string) (data.Session, error) {
	m.GetSessionCalls = append(
		m.GetSessionCalls,
		GetSessionCall{
			Token: token,
		},
	)
	return m.GetSessionResponse.Session, m.GetSessionResponse.Error
}

type ExtendSessionCall struct {
	Token string
}

type ExtendSessionResponse struct {
	Error error
}

func (m *MockPostgres) SetExtendSessionResponse(response ExtendSessionResponse) {
	m.ExtendSessionResponse = response
}

func (m *MockPostgres) ExtendSession(token string) error {
	m.ExtendSessionCalls = append(
		m.ExtendSessionCalls, ExtendSessionCall{
			Token: token,
		},
	)
	return m.ExtendSessionResponse.Error
}
