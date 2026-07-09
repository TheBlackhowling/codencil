package models

import (
	"time"

	"github.com/TheBlackHowling/typrow"
)

// User is an authenticated account (dev external_id or OIDC subject).
type User struct {
	typrow.Model
	ID          string    `db:"id" load:"primary"`
	OrgID       string    `db:"org_id"`
	ExternalID  string    `db:"external_id"`
	DisplayName string    `db:"display_name"`
	CreatedAt   time.Time `db:"created_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) QueryByID() string {
	return `SELECT id, org_id, external_id, display_name, created_at FROM users WHERE id = $1`
}

func init() {
	typrow.RegisterModel[*User]()
}
