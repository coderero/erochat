package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/coderero/erochat-server/db/mysql/queries"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	driver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// ProfileStore is a MySQL data store for profile.
type ProfileStore struct {
	// ConnectionPool is a pool of connections to the database.
	pool *ConnectionPool
}

// NewProfileStore creates a new ProfileStore.
func NewProfileStore(pool *ConnectionPool) *ProfileStore {
	return &ProfileStore{
		pool: pool,
	}
}

// GetByID returns a profile by its uuid.
func (s *ProfileStore) GetByUID(id uuid.UUID) (*types.Profile, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	profile := &types.Profile{}
	err = db.QueryRow(queries.GetUserProfileByUID, id).Scan(&profile.ID, &profile.UID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Bio, &profile.Avatar, &profile.Username, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt)
	if err != nil {
		// If the profile is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrProfileNotFound
		}
		return nil, interfaces.ErrFailedToGetProfile
	}
	return profile, nil
}

// GetByUserID returns a profile by its user id.
func (s *ProfileStore) GetByUserID(id int) (*types.Profile, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	profile := &types.Profile{}
	err = db.QueryRow(queries.GetUserProfileByUserID, id).Scan(&profile.ID, &profile.UID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Bio, &profile.Avatar, &profile.Username, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt)
	if err != nil {
		// If the profile is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrProfileNotFound
		}
		return nil, interfaces.ErrFailedToGetProfile
	}
	return profile, nil
}

// GetByEmail returns a profile by its email.
func (s *ProfileStore) GetByEmail(email string) (*types.Profile, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	profile := &types.Profile{}
	err = db.QueryRow(queries.GetUserProfileByEmail, email).Scan(&profile.ID, &profile.UID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Bio, &profile.Avatar, &profile.Username, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt)
	if err != nil {
		// If the profile is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrProfileNotFound
		}
		return nil, interfaces.ErrFailedToGetProfile
	}
	return profile, nil
}

// Create creates a new profile.
func (s *ProfileStore) Create(profile *types.Profile) (*types.Profile, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	result, err := db.Exec(queries.CreateUserProfile, profile.UID, profile.UserID, profile.FirstName, profile.LastName, profile.Bio, profile.Avatar)
	if err != nil {
		// Check if the profile already exists.
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, interfaces.ErrProfileExists
		}
		return nil, interfaces.ErrFailedToCreateProfile
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, interfaces.ErrFailedToCreateProfile
	}

	var uid string
	row := db.QueryRow(queries.GetUserProfileByID, id)
	err = row.Scan(&profile.ID, &uid, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Bio, &profile.Avatar, &profile.Username, &profile.Email, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt)
	if err != nil {
		return nil, interfaces.ErrFailedToCreateProfile
	}

	profile.UID, err = uuid.Parse(uid)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfileStore) CreateFriendship(userID, friendID string) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	a, err := db.Exec(queries.CreateFriendship, userID, friendID)
	if err != nil {
		sqlErr, _ := err.(*driver.MySQLError)
		if string(sqlErr.SQLState[:]) == "45000" && strings.Contains(sqlErr.Message, "duplicate request") {

			return interfaces.ErrDuplicateFriendship
		}
		if string(sqlErr.SQLState[:]) == "45000" && strings.Contains(sqlErr.Message, "self request") {
			return interfaces.ErrSelfFriendship
		}
		return interfaces.ErrFailedToCreateFriendship
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return interfaces.ErrFriendNotFound
	}

	return nil
}

// Update updates a profile.
func (s *ProfileStore) Update(profile *types.Profile) (*types.Profile, error) {
	var p types.Profile
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	_, err = db.Exec(queries.UpdateUserProfile, profile.FirstName, profile.LastName, profile.Bio, profile.Avatar, profile.UserID)
	// SetStatus sets the status of a user.
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, interfaces.ErrFailedToUpdateProfile
	}

	row := db.QueryRow(queries.GetUserProfileByUserID, profile.UserID)
	err = row.Scan(&p.ID, &p.UID, &p.UserID, &p.FirstName, &p.LastName, &p.Bio, &p.Avatar, &p.Username, &p.Email, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		fmt.Printf("error3: %v\n", err)
		return nil, interfaces.ErrFailedToUpdateProfile
	}

	return &p, nil
}

// Delete deletes a profile by its uuid.
func (s *ProfileStore) Delete(id int) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	_, err = db.Exec(queries.DeleteUserProfile, id)
	if err != nil {
		return interfaces.ErrFailedToDeleteProfile
	}

	return nil
}

// Reactivate reactivates a profile by its uuid.
func (s *ProfileStore) Reactivate(id int) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	_, err = db.Exec(queries.ReactivateUserProfile, id)
	if err != nil {
		return interfaces.ErrFailedToUpdateProfile
	}

	return nil
}
