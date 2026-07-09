-- Users (Phase 4 auth).

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id TEXT NOT NULL DEFAULT 'dev',
    external_id TEXT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, external_id)
);

CREATE INDEX idx_users_org_external ON users (org_id, external_id);
