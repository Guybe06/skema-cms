CREATE TABLE collections (
    id              UUID PRIMARY KEY,
    connection_id   UUID NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    table_name      VARCHAR(255) NOT NULL,
    display_name    VARCHAR(255),
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(connection_id, table_name)
);

CREATE TABLE collection_fields (
    id              UUID PRIMARY KEY,
    collection_id   UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    column_name     VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    required        BOOLEAN NOT NULL DEFAULT FALSE,
    is_unique       BOOLEAN NOT NULL DEFAULT FALSE,
    default_value   TEXT,
    options         JSONB,
    position        INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(collection_id, column_name)
);

CREATE INDEX idx_collections_org ON collections(organization_id);
CREATE INDEX idx_collections_connection ON collections(connection_id);
CREATE INDEX idx_collection_fields_collection ON collection_fields(collection_id);
