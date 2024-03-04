package interfaces

import (
	"errors"

	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrFailedToGetUser is returned when the user is not found.
	ErrFailedToGetUser = errors.New("failed to get user")

	// ErrFailedToCreateUser is returned when the user is not found.
	ErrFailedToCreateUser = errors.New("failed to create user")

	// ErrFailedToUpdateUser is returned when the user is not found.
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// ErrFailedToDeleteUser is returned when the user is not found.
	ErrFailedToDeleteUser = errors.New("failed to delete user")

	// ErrEmailExists is returned when the email exists.
	ErrEmailExists = errors.New("email exists")

	// ErrUsernameExists is returned when the username exists.
	ErrUsernameExists = errors.New("username exists")
)

type UserStore interface {
	// GetByID returns a user by its uuid.
	GetByID(uuid uuid.UUID) (*types.User, error)

	// GetByEmail returns a user by its email.
	GetByEmail(email string) (*types.User, error)

	// GetByUsername returns a user by its username.
	GetByUsername(username string) (*types.User, error)

	// Create creates a new user.
	Create(user *types.User) (*types.User, error)

	// Update updates a user.
	Update(id uuid.UUID, user *types.User) (*types.User, error)

	// Delete deletes a user by its uuid.
	Delete(id uuid.UUID) (uuid.UUID, error)
}
