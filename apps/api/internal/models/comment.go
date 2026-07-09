package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

// Comment is one message in a comment thread.
type Comment struct {
	typrow.Model
	ID        string    `db:"id" load:"primary"`
	ThreadID  string    `db:"thread_id"`
	AuthorID  string    `db:"author_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

func (c *Comment) TableName() string {
	return "comments"
}

func (c *Comment) QueryByID() string {
	return `SELECT id, thread_id, author_id, body, created_at FROM comments WHERE id = $1`
}

func init() {
	typrow.RegisterModel[*Comment]()
}
