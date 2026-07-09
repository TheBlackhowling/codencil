package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

// Document is the editable draft and metadata for a markdown document.
type Document struct {
	typrow.Model
	ID            string    `db:"id" load:"primary"`
	OrgID         string    `db:"org_id"`
	Title         string    `db:"title"`
	DraftMarkdown string    `db:"draft_markdown"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (d *Document) TableName() string {
	return "documents"
}

func (d *Document) QueryByID() string {
	return `SELECT id, org_id, title, draft_markdown, created_at, updated_at
		FROM documents WHERE id = $1`
}

func init() {
	typrow.RegisterModel[*Document]()
}
