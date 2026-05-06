CREATE TABLE api_keys (
    id                  UUID PRIMARY KEY,
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    key_hash            TEXT UNIQUE NOT NULL,
    key_prefix          VARCHAR(12) NOT NULL,
    permissions         JSONB NOT NULL DEFAULT '{"read":true,"create":false,"update":false,"delete":false}',
    allowed_collections JSONB,
    expires_at          TIMESTAMPTZ,
    last_used_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_org ON api_keys(organization_id);
