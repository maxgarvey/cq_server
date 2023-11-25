CREATE TABLE sessions IF NOT EXISTS (
    user_id INT NOT NULL,
    token VARCHAR(192) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    good_until TIMESTAMP NOT NULL
);