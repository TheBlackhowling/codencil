-- Review path: comment threads, version anchors, and comments (Phase 2).

CREATE TABLE comment_threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_comment_threads_document_id ON comment_threads (document_id);

CREATE TABLE version_anchors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anchor_id UUID NOT NULL,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version INT NOT NULL CHECK (version > 0),
    thread_id UUID NOT NULL REFERENCES comment_threads(id) ON DELETE CASCADE,
    start_line INT NOT NULL CHECK (start_line > 0),
    end_line INT NOT NULL CHECK (end_line >= start_line),
    quoted_text TEXT NOT NULL,
    anchor_status TEXT NOT NULL DEFAULT 'active'
        CHECK (anchor_status IN ('active', 'shifted', 'orphaned')),
    review_state TEXT NOT NULL DEFAULT 'open'
        CHECK (review_state IN ('open', 'resolved')),
    resolved_by TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id, version) REFERENCES document_versions (document_id, version) ON DELETE CASCADE,
    UNIQUE (document_id, version, anchor_id)
);

CREATE INDEX idx_version_anchors_document_version ON version_anchors (document_id, version);
CREATE INDEX idx_version_anchors_thread_id ON version_anchors (thread_id);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES comment_threads(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_comments_thread_id ON comments (thread_id);
