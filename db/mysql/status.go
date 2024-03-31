package mysql

import (
	"fmt"

	"github.com/coderero/erochat-server/db/mysql/queries"
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

// StatusStore is a MySQL data store for status.
type StatusStore struct {
	// ConnectionPool is a pool of connections to the database.
	pool *ConnectionPool
}

// NewStatusStore creates a new StatusStore.
func NewStatusStore(pool *ConnectionPool) *StatusStore {
	return &StatusStore{
		pool: pool,
	}
}

// GetStatus gets the status of a user.
func (s *StatusStore) GetStatus(id uuid.UUID) ([]*types.UserStatus, error) {
	var statuses []*types.UserStatus
	statuses = []*types.UserStatus{}
	db, err := s.pool.Get()
	if err != nil {
		return statuses, err
	}
	defer s.pool.Release()

	rows, err := db.Query(queries.GetUsersStatus, id)
	if err != nil {
		return statuses, err
	}
	defer rows.Close()

	for rows.Next() {
		status := &types.UserStatus{}
		err = rows.Scan(&status.ID, &status.UID, &status.UserID, &status.Title, &status.ResourceURI, &status.ResourceThumbnail, &status.CreatedAt)
		if err != nil {
			return statuses, err
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (s *StatusStore) GetStatusByUID(userUID, uid uuid.UUID) (*types.UserStatus, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	status := &types.UserStatus{}
	err = db.QueryRow(queries.GetStatusByID, userUID, uid).Scan(&status.ID, &status.UID, &status.UserID, &status.Title, &status.ResourceURI, &status.ResourceThumbnail, &status.CreatedAt)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// CreateStatus creates a new status.
func (s *StatusStore) CreateStatus(status *types.UserStatus) (*types.UserStatus, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	st, err := db.Exec(queries.CreateStatus, status.UserID, status.Title, status.ResourceURI, status.ResourceThumbnail)
	if err != nil {
		return nil, err
	}

	id, err := st.LastInsertId()
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(queries.GetStatusByID, id).Scan(&status.ID, &status.UID, &status.UserID, &status.Title, &status.ResourceURI, &status.ResourceThumbnail, &status.CreatedAt)
	if err != nil {
		return nil, err
	}
	return status, nil

}

// DeleteStatus deletes a status by its id.
func (s *StatusStore) DeleteStatus(userID uuid.UUID, uid uuid.UUID) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	a, err := db.Exec(queries.DeleteStatus, userID, uid)
	if err != nil {
		return err
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return fmt.Errorf("status not found")
	}

	return nil
}
