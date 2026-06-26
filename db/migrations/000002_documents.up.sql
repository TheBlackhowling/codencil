-- Documents and published versions (Phase 1 read path).

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    draft_markdown TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_documents_org_id ON documents (org_id);

CREATE TABLE document_versions (
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INT NOT NULL CHECK (version > 0),
    markdown TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_by TEXT NOT NULL,
    PRIMARY KEY (document_id, version)
);

CREATE INDEX idx_document_versions_document_id ON document_versions (document_id);
