package admin

import (
	"fmt"
)

func (a *Admin) ValidateSession(token string) (bool, error) {
	session, err := a.Postgres.GetSession(token)
	if err != nil {
		a.Logger.Error(
			fmt.Sprintf(
				"error looking up session: %s\n",
				fmt.Errorf("%w", err),
			),
		)
		return false, err
	}

	if session.GoodUntil.Before(a.Clock.Now()) {
		return false, nil
	}
	return true, nil
}

func (a *Admin) ExtendSession(token string) error {
	return a.Postgres.ExtendSession(token)
}
