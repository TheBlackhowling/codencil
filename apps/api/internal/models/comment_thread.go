package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

// CommentThread groups comments on a document across versions.
type CommentThread struct {
	typrow.Model
	ID         string    `db:"id" load:"primary"`
	DocumentID string    `db:"document_id"`
	CreatedAt  time.Time `db:"created_at"`
}

func (t *CommentThread) TableName() string {
	return "comment_threads"
}

func (t *CommentThread) QueryByID() string {
	return `SELECT id, document_id, created_at FROM comment_threads WHERE id = $1`
}

func init() {
	typrow.RegisterModel[*CommentThread]()
}
