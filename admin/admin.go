package admin

import (
	"errors"
	"fmt"

	"github.com/maxgarvey/cq_server/postgres"
)

type Admin struct {
	Postgres *postgres.Postgres
}

func (a *Admin) Login(username string, password string) (bool, error) {
	// Check if the user exists.
	exists, err := a.Postgres.UserExists(username, password)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New(
			fmt.Sprintf(
				"No user found for username=%s",
				username,
			),
		)
	}

	// Update the logged in time.
	err = a.Postgres.UpdateLastLogin(username)
	if err != nil {
		return true, err
	}

	return true, nil
}
