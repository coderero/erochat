package interfaces

import (
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

type StatusStore interface {
	// GetStatus gets the status of a user.
	GetStatus(uid uuid.UUID) ([]*types.UserStatus, error)

	// GetStatusByUID gets a status by its id.
	GetStatusByUID(userUID, uid uuid.UUID) (*types.UserStatus, error)

	// CreateStatus creates a new status.
	CreateStatus(status *types.UserStatus) (*types.UserStatus, error)

	// DeleteStatus deletes a status by its id.
	DeleteStatus(userID uuid.UUID, uid uuid.UUID) error
}
