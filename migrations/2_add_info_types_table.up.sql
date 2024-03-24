CREATE TABLE IF NOT EXISTS info_types (
    info_id SERIAL PRIMARY KEY,
    typename VARCHAR(32) NOT NULL CHECK (LENGTH(info_type) >= 1),
    owner_fk BIGINT DEFAULT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX owner_fk_idx ON info_types (owner_fk NULLS FIRST);
CREATE INDEX typename_idx ON info_types (typename);