CREATE TABLE IF NOT EXISTS users_info (
    user_fk BIGINT REFERENCES users(user_id) NOT NULL ON DELETE CASCADE,
    info_type_fk INT REFERENCES info_types(info_id) DEFAULT NULL ON DELETE SET NULL,
    service_fk BIGINT REFERENCES services(service_id) ON DELETE CASCADE,
    info VARCHAR(1024) NOT NULL,
    UNIQUE(info_type_fk, service_fk, info)
);