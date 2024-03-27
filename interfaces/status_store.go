package interfaces

import (
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

type StatusStore interface {
	// GetStatus gets the status of a user.
	GetStatus(id int) ([]*types.UserStatus, error)

	// CreateStatus creates a new status.
	CreateStatus(status *types.UserStatus) (*types.UserStatus, error)

	// DeleteStatus deletes a status by its id.
	DeleteStatus(userID int, id uuid.UUID) error
}
