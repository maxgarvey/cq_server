CREATE TABLE IF NOT EXISTS cq_server_users (
	user_id serial PRIMARY KEY,
	username varchar(192) UNIQUE NOT NULL,
	password varchar(192) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	last_login TIMESTAMP
);
