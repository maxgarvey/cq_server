package admin

type MockAdmin struct {
	ExtendSessionCalls      []ExtendSessionCall
	ExtendSessionResponse   ExtendSessionResponse
	LoginCalls              []LoginCall
	LoginResponse           LoginResponse
	ValidateSessionCalls    []ValidateSessionCall
	ValidateSessionResponse ValidateSessionResponse
}

func InitMock() MockAdmin {
	return MockAdmin{
		ExtendSessionCalls: []ExtendSessionCall{},
		ExtendSessionResponse: ExtendSessionResponse{
			Error: nil,
		},
		LoginCalls: []LoginCall{},
		LoginResponse: LoginResponse{
			Token: "session",
			Error: nil,
		},
		ValidateSessionCalls: []ValidateSessionCall{},
		ValidateSessionResponse: ValidateSessionResponse{
			Valid: true,
			Error: nil,
		},
	}
}

type ExtendSessionCall struct {
	Token string
}

type ExtendSessionResponse struct {
	Error error
}

func (m *MockAdmin) SetExtendSessionResponse(response ExtendSessionResponse) {
	m.ExtendSessionResponse = response
}

func (m *MockAdmin) ExtendSession(token string) error {
	m.ExtendSessionCalls = append(
		m.ExtendSessionCalls,
		ExtendSessionCall{Token: token},
	)
	return m.ExtendSessionResponse.Error
}

type LoginCall struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
	Error error
}

func (m *MockAdmin) SetLoginResponse(response LoginResponse) {
	m.LoginResponse = response
}

func (m *MockAdmin) Login(username string, password string) (string, error) {
	m.LoginCalls = append(
		m.LoginCalls,
		LoginCall{
			Username: username,
			Password: password,
		},
	)
	return m.LoginResponse.Token, m.LoginResponse.Error
}

type ValidateSessionCall struct {
	Token string
}

type ValidateSessionResponse struct {
	Valid bool
	Error error
}

func (m *MockAdmin) SetValidateSessionResponse(response ValidateSessionResponse) {
	m.ValidateSessionResponse = response
}

func (m *MockAdmin) ValidateSession(token string) (bool, error) {
	m.ValidateSessionCalls = append(
		m.ValidateSessionCalls,
		ValidateSessionCall{
			Token: token,
		},
	)
	return m.ValidateSessionResponse.Valid, m.ValidateSessionResponse.Error
}
