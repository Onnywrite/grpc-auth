CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    login VARCHAR(30) UNIQUE NOT NULL CHECK (LENGTH(login) >= 1),
    email VARCHAR(255) UNIQUE NULLS NOT DISTINCT,
    phone VARCHAR(16) UNIQUE NULLS NOT DISTINCT,
    password VARCHAR(1024) NOT NULL,
    deleted_at DATE DEFAULT NULL
);