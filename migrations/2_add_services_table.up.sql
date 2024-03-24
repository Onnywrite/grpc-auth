CREATE TABLE IF NOT EXISTS services (
    service_id BIGSERIAL PRIMARY KEY,
    owner_fk BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(32) NOT NULL CHECK (LENGTH(name) >= 1),
    deleted_at DATE DEFAULT NULL,
    UNIQUE(owner_fk, name)
);