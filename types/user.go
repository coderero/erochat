package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents a user.
type User struct {
	ID        int          `json:"id" db:"id"`
	UID       uuid.UUID    `json:"uid" db:"uid"`
	Username  string       `json:"username" db:"username"`
	Email     string       `json:"email" db:"email"`
	Password  string       `json:"password" db:"password"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}
