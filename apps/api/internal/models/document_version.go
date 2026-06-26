package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

// DocumentVersion is an immutable published markdown snapshot.
type DocumentVersion struct {
	typrow.Model
	DocumentID  string    `db:"document_id" load:"composite:document_version"`
	Version     int       `db:"version" load:"composite:document_version"`
	Markdown    string    `db:"markdown"`
	PublishedAt time.Time `db:"published_at"`
	PublishedBy string    `db:"published_by"`
}

func (v *DocumentVersion) TableName() string {
	return "document_versions"
}

// QueryByDocumentIDVersion matches composite key field order (DocumentID, Version).
func (v *DocumentVersion) QueryByDocumentIDVersion() string {
	return `SELECT document_id, version, markdown, published_at, published_by
		FROM document_versions WHERE document_id = $1 AND version = $2`
}

func init() {
	typrow.RegisterModel[*DocumentVersion]()
}
