package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

const (
	RoleOwner    = "owner"
	RoleReviewer = "reviewer"
	RoleViewer   = "viewer"
)

// DocumentMember links a user to a document with a role.
type DocumentMember struct {
	typrow.Model
	DocumentID string    `db:"document_id" load:"composite:document_member"`
	UserID     string    `db:"user_id" load:"composite:document_member"`
	Role       string    `db:"role"`
	CreatedAt  time.Time `db:"created_at"`
}

func (m *DocumentMember) TableName() string {
	return "document_members"
}

func (m *DocumentMember) QueryByDocumentIDUserID() string {
	return `SELECT document_id, user_id, role, created_at
		FROM document_members WHERE document_id = $1 AND user_id = $2`
}

func init() {
	typrow.RegisterModel[*DocumentMember]()
}
