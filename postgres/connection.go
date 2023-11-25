package postgres

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/benbjohnson/clock"
	_ "github.com/lib/pq"
	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/data"
	"golang.org/x/crypto/bcrypt"
)

type Postgres struct {
	Clock      clock.Clock
	Connection *sql.DB
	Logger     *slog.Logger
}

func ConfigInit(config config.Postgres, clock clock.Clock, logger *slog.Logger) Postgres {
	return Init(config.Host, config.Port, config.Username, config.Password, config.DBName, logger, clock)
}

func Init(host string, port int, username string, password string, dbname string, logger *slog.Logger, clock clock.Clock) Postgres {
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
		Clock:      clock,
		Connection: connection,
		Logger:     logger,
	}
}

func (p *Postgres) GetUser(username string, password string) (data.User, error) {
	var user data.User
	if err := p.Connection.QueryRow(
		"SELECT user_id, username, created_at, last_login "+
			"FROM cq_server_users "+
			"WHERE username=$1 "+
			"AND password=$2",
		username,
		password,
	).Scan(
		&user.ID, &user.Username, &user.CreatedAt, &user.LastLogin,
	); err != nil {
		return user, err
	}

	return user, nil
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

func GenerateToken(user_id int, timestamp string) (string, error) {
	tokenString := fmt.Sprintf("%d%s", user_id, timestamp)
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(tokenString),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func (p *Postgres) CreateSession(user_id int) (string, error) {
	session_token, err := GenerateToken(
		user_id,
		p.Clock.Now().String(),
	)
	if err != nil {
		return "", err
	}

	if _, err := p.Connection.Exec(
		"INSERT INTO sessions "+
			"(user_id, token, created_at, good_until) "+
			"VALUES "+
			"($1, $2, $3, $4)",
		user_id,
		session_token,
		p.Clock.Now().String(),
		p.Clock.Now().Add(time.Hour*24).String(),
	); err != nil {
		return "", err
	}

	return session_token, nil
}
