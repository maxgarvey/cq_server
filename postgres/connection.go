package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Connection *sql.DB
}

func Init(host string, port int, username string, password string, dbname string, logger slog.Logger) Postgres {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname,
	)

	connection, err := sql.Open(
		"postgres",
		connectionString,
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error connecting to postgres: %s",
				err.Error(),
			),
		)
		panic(err)
	}

	return Postgres{
		Connection: connection,
	}
}

func (p *Postgres) UserExists(username string, password string) (bool, error) {
	var count int
	if err := p.Connection.QueryRow(
		"SELECT COUNT(*) "+
			"FROM cq_server_users "+
			"WHERE username=$1 "+
			"AND password=$2",
		username,
		password,
	).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *Postgres) UpdateLastLogin(username string) error {
	if _, err := p.Connection.Exec(
		"UPDATE cq_server_users "+
			"SET last_login=NOW() "+
			"WHERE username=$1",
		username,
	); err != nil {
		return err
	}

	return nil
}
