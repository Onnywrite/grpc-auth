CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    login VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NULLS DISTINCT,
    phone VARCHAR(16) UNIQUE NULLS DISTINCT,
    password VARCHAR(1024) NOT NULL
);

CREATE UNIQUE INDEX login_idx ON users (login);
CREATE INDEX email_idx ON users (email);
CREATE INDEX phone_idx ON users (phone);