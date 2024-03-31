package types

import "github.com/google/uuid"

// UserStatus is the model for the user status.
type UserStatus struct {
	// ID is the unique identifier of the user status.
	ID                int       `json:"-"`
	UID               uuid.UUID `json:"uid"`
	UserID            uuid.UUID `json:"-"`
	Title             string    `json:"title"`
	ResourceURI       string    `json:"resource_uri"`
	ResourceThumbnail string    `json:"resource_thumbnail"`
	CreatedAt         string    `json:"-"`
}
