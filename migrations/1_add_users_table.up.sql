CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    nickname VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NULLS DISTINCT,
    phone VARCHAR(16) UNIQUE NULLS DISTINCT,
    password VARCHAR(60) NOT NULL,
    deleted_at DATE DEFAULT NULL
);

CREATE UNIQUE INDEX nickname_idx ON users (nickname);
CREATE INDEX email_idx ON users (email);
CREATE INDEX phone_idx ON users (phone);