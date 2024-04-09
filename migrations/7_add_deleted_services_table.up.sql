CREATE TABLE deleted_services (
    service_fk BIGINT PRIMARY KEY REFERENCES services(service_id) ON DELETE CASCADE,
    deleted_at DATE NOT NULL
);