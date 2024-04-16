CREATE TABLE IF NOT EXISTS sessions (
    session_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_fk BIGINT NOT NULL,
    service_fk BIGINT NOT NULL,
    ip CIDR NOT NULL,
    browser VARCHAR(32),
    os VARCHAR(16),
    at TIMESTAMP NOT NULL DEFAULT NOW(),
    terminated_at TIMESTAMP DEFAULT NULL,
    UNIQUE(user_fk, service_fk, ip, browser, os),
    FOREIGN KEY (user_fk, service_fk) REFERENCES signups(user_fk, service_fk) ON DELETE CASCADE
);

CREATE INDEX sessions_fk_idx ON sessions (user_fk, service_fk);
CREATE INDEX sessions_info_idx ON sessions (ip, browser, os);