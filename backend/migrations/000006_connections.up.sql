CREATE TABLE connections (
    id                  UUID PRIMARY KEY,
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    driver              VARCHAR(50) NOT NULL,
    host_encrypted      TEXT NOT NULL,
    port_encrypted      TEXT NOT NULL,
    database_encrypted  TEXT NOT NULL,
    user_encrypted      TEXT NOT NULL,
    password_encrypted  TEXT NOT NULL,
    ssl_mode            VARCHAR(50) NOT NULL DEFAULT 'disable',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_connections_org ON connections(organization_id);
