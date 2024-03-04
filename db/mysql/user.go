package mysql

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/coderero/erochat-server/db/mysql/queries"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

// UserStore is a MySQL data store for user.
type UserStore struct {
	// ConnectionPool is a pool of connections to the database.
	pool *ConnectionPool
}

// CheckUserExistsResult represents the result of checking if a user exists.
type CheckUserExistsResult struct {
	UUIDExists     bool
	EmailExists    bool
	UsernameExists bool
}

// NewUserStore creates a new UserStore.
func NewUserStore(pool *ConnectionPool) *UserStore {
	return &UserStore{
		pool: pool,
	}
}

// GetByID returns a user by its uuid.
func (s *UserStore) GetByID(id uuid.UUID) (*types.User, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	err = db.QueryRow(queries.GetUserByID, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, interfaces.ErrFailedToGetUser
	}
	return user, nil
}

// GetByEmail returns a user by its email.
func (s *UserStore) GetByEmail(email string) (*types.User, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	err = db.QueryRow(queries.GetUserByEmail, email).Scan(&user.ID, &user.UID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, interfaces.ErrFailedToGetUser
	}
	return user, nil
}

// GetByUsername returns a user by its username.
func (s *UserStore) GetByUsername(username string) (*types.User, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	err = db.QueryRow(queries.GetUserByUsername, username).Scan(&user.ID, &user.UID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, interfaces.ErrFailedToGetUser
	}
	return user, nil
}

// Create creates a new user.
func (s *UserStore) Create(user *types.User) (*types.User, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec(queries.CreateUser, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, checkForErrorConstraint(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, interfaces.ErrFailedToCreateUser
	}

	err = db.QueryRow(queries.GetUserByID, id).Scan(&user.ID, &user.UID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Printf("error: %v", err)
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, interfaces.ErrFailedToCreateUser
	}
	return user, nil
}

// Update updates a user.
func (s *UserStore) Update(id uuid.UUID, user *types.User) (*types.User, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	// Update a user.
	row := db.QueryRow(queries.UpdateUser, user.Username, user.Email, user.Password, user.UpdatedAt, id, id)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, checkForErrorConstraint(err)
		}
		return nil, interfaces.ErrFailedToUpdateUser
	}
	return user, nil
}

// Delete deletes a user by its uuid.
func (s *UserStore) Delete(id uuid.UUID) (uuid.UUID, error) {
	db, err := s.pool.Get()
	if err != nil {
		return uuid.Nil, err
	}

	// Delete a user.
	row := db.QueryRow(queries.DeleteUser, id, id)
	var deletedID uuid.UUID
	err = row.Scan(&deletedID)
	if err != nil {
		// If the user is not found, return an error.
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, interfaces.ErrUserNotFound
		}
		return uuid.Nil, interfaces.ErrFailedToDeleteUser
	}
	return deletedID, nil
}

func checkForErrorConstraint(err error) error {
	if strings.Contains(err.Error(), "email") {
		return interfaces.ErrEmailExists
	} else if strings.Contains(err.Error(), "username") {
		return interfaces.ErrUsernameExists
	}
	return interfaces.ErrFailedToGetUser
}
