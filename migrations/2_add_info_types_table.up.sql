CREATE TABLE IF NOT EXISTS info_types (
    info_id SERIAL PRIMARY KEY,
    typename VARCHAR(32) NOT NULL CHECK (LENGTH(info_type) >= 1),
    owner_fk BIGINT DEFAULT NULL REFERENCES users(user_id) ON DELETE CASCADE
);