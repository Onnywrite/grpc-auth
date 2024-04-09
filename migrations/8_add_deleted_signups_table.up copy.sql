CREATE TABLE deleted_signups (
    signup_fk BIGINT PRIMARY KEY REFERENCES signups(signup_id) ON DELETE CASCADE,
    deleted_at DATE NOT NULL
);