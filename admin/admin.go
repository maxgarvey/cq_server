package admin

import (
	"fmt"
	"log/slog"

	"github.com/maxgarvey/cq_server/postgres"
)

type Admin struct {
	Postgres *postgres.Postgres
	Logger   *slog.Logger
}

func (a *Admin) Login(username string, password string) (string, error) {
	user, err := a.Postgres.GetUser(
		username,
		password,
	)
	if err != nil {
		a.Logger.Error(
			fmt.Sprintf(
				"error looking up user: %s\n",
				fmt.Errorf("%w", err),
			),
		)
		return "", err
	}

	// TODO: figure out what happens for empty user
	if user.ID < 0 {
		return "", fmt.Errorf(
			"invalid user id: %d",
			user.ID,
		)
	}
	token, err := a.Postgres.CreateSession(user.ID)
	if err != nil {
		a.Logger.Error(
			fmt.Sprintf(
				"error looking up user: %s\n",
				fmt.Errorf("%w", err),
			),
		)
		return "", err
	}

	err = a.Postgres.UpdateLastLogin(user.Username)
	if err != nil {
		a.Logger.Error(
			fmt.Sprintf(
				"error updating last login: %s\n",
				fmt.Errorf("%w", err),
			),
		)
		return "", err
	}
	return token, nil
}
