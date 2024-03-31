package types

import (
	"time"

	"github.com/google/uuid"
)

// Friend is struct for friendship join table.
type Friend struct {
	RID        uuid.UUID  `json:"rid"`
	UID        uuid.UUID  `json:"uid"`
	Username   string     `json:"username"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"lastname"`
	Bio        string     `json:"bio"`
	Avatar     string     `json:"avatar"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
}

type FriendStatus struct {
	RID               uuid.UUID `json:"rid"`
	UID               uuid.UUID `json:"uid"`
	StatusID          uuid.UUID `json:"status_id"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Avatar            string    `json:"avatar"`
	Title             string    `json:"title"`
	ResourceURI       string    `json:"resource_uri"`
	ResourceThumbnail string    `json:"resource_thumbnail"`
}
