CREATE TABLE deleted_users (
    user_fk BIGINT PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    deleted_at DATE NOT NULL
);