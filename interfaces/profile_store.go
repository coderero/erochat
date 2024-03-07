package interfaces

import (
	"errors"

	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

var (
	// ErrProfileNotFound is returned when the profile is not found.
	ErrProfileNotFound = errors.New("profile not found")

	// ErrFailedToGetProfile is returned when the profile is not found.
	ErrFailedToGetProfile = errors.New("failed to get profile")

	// ErrFailedToCreateProfile is returned when the profile is not found.
	ErrFailedToCreateProfile = errors.New("failed to create profile")

	// ErrFailedToUpdateProfile is returned when the profile is not found.
	ErrFailedToUpdateProfile = errors.New("failed to update profile")

	// ErrFailedToDeleteProfile is returned when the profile is not found.
	ErrFailedToDeleteProfile = errors.New("failed to delete profile")

	// ErrProfileExists is returned when the profile exists.
	ErrProfileExists = errors.New("profile exists")
)

// ProfileStore is a data store for profile.
type ProfileStore interface {
	// GetByUID returns a profile by its uuid.
	GetByUID(id uuid.UUID) (*types.Profile, error)

	// GetByUserID returns a profile by its user id.
	GetByUserID(id int) (*types.Profile, error)

	// GetByEmail returns a profile by its email.
	GetByEmail(email string) (*types.Profile, error)

	// Create creates a new profile.
	Create(profile *types.Profile) (*types.Profile, error)

	// Update updates a profile.
	Update(profile *types.Profile) (*types.Profile, error)

	// Delete deletes a profile by its uuid.
	Delete(id int) (int, error)
}