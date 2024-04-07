CREATE TABLE IF NOT EXISTS sessions (
    session_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    signup_fk BIGINT REFERENCES signups(signup_id) ON DELETE CASCADE,
    ip CIDR NOT NULL,
    browser VARCHAR(32),
    os VARCHAR(16),
    at TIMESTAMP NOT NULL DEFAULT NOW(),
    terminated_at TIMESTAMP DEFAULT NULL,
    UNIQUE(signup_fk, ip, browser, os)
);

CREATE INDEX signup_fk_idx ON sessions (signup_fk);