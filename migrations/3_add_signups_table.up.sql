CREATE TABLE IF NOT EXISTS signups (
    user_fk BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    service_fk BIGINT REFERENCES services(service_id) ON DELETE CASCADE,
    at TIMESTAMP NOT NULL,
    banned_at TIMESTAMP DEFAULT NULL,
    deleted_at DATE DEFAULT NULL,
    PRIMARY KEY(user_fk, service_fk)
);