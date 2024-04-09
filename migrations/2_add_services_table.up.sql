CREATE TABLE IF NOT EXISTS services (
    service_id BIGSERIAL PRIMARY KEY,
    owner_fk BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(32) NOT NULL,
    UNIQUE(owner_fk, name)
);

CREATE INDEX owner_fk_idx ON services (owner_fk);