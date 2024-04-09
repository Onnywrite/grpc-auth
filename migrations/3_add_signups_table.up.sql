CREATE TABLE IF NOT EXISTS signups (
    signup_id BIGSERIAL PRIMARY KEY,
    user_fk BIGINT REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    service_fk BIGINT REFERENCES services(service_id) ON DELETE CASCADE NOT NULL,
    at TIMESTAMP NOT NULL DEFAULT NOW(),
    banned_at TIMESTAMP DEFAULT NULL,
    deleted_at DATE DEFAULT NULL,
    UNIQUE(user_fk, service_fk)
);