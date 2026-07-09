package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

const (
	AnchorStatusActive   = "active"
	AnchorStatusShifted  = "shifted"
	AnchorStatusOrphaned = "orphaned"

	ReviewStateOpen     = "open"
	ReviewStateResolved = "resolved"
)

// VersionAnchor ties a comment thread to quoted text on a published version.
type VersionAnchor struct {
	typrow.Model
	ID           string     `db:"id" load:"primary"`
	AnchorID     string     `db:"anchor_id"`
	DocumentID   string     `db:"document_id"`
	Version      int        `db:"version"`
	ThreadID     string     `db:"thread_id"`
	StartLine    int        `db:"start_line"`
	EndLine      int        `db:"end_line"`
	QuotedText   string     `db:"quoted_text"`
	AnchorStatus string     `db:"anchor_status"`
	ReviewState  string     `db:"review_state"`
	ResolvedBy   *string    `db:"resolved_by"`
	ResolvedAt   *time.Time `db:"resolved_at"`
	CreatedAt    time.Time  `db:"created_at"`
}

func (a *VersionAnchor) TableName() string {
	return "version_anchors"
}

func (a *VersionAnchor) QueryByID() string {
	return `SELECT id, anchor_id, document_id, version, thread_id, start_line, end_line,
		quoted_text, anchor_status, review_state, resolved_by, resolved_at, created_at
		FROM version_anchors WHERE id = $1`
}

func init() {
	typrow.RegisterModel[*VersionAnchor]()
}
