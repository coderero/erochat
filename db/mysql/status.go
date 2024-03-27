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
func (s *StatusStore) GetStatus(id int) ([]*types.UserStatus, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(queries.GetUsersStatus, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*types.UserStatus
	for rows.Next() {
		status := &types.UserStatus{}
		err = rows.Scan(&status.ID, &status.UID, &status.UserID, &status.Title, &status.ResourceURI, &status.ResourceThumbnail, &status.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// CreateStatus creates a new status.
func (s *StatusStore) CreateStatus(status *types.UserStatus) (*types.UserStatus, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}

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
		fmt.Println(err)
		return nil, err
	}
	return status, nil

}

// DeleteStatus deletes a status by its id.
func (s *StatusStore) DeleteStatus(userID int, id uuid.UUID) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}

	a, err := db.Exec(queries.DeleteStatus, id, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return fmt.Errorf("status not found")
	}

	return nil
}
