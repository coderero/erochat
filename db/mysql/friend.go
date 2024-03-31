package mysql

import (
	"database/sql"
	"errors"

	"github.com/coderero/erochat-server/db/mysql/queries"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

// FriendStore is a MySQL data store for friend.
type FriendStore struct {
	// ConnectionPool is a pool of connections to the database.
	pool *ConnectionPool
}

// NewFriendStore creates a new FriendStore.
func NewFriendStore(pool *ConnectionPool) *FriendStore {
	return &FriendStore{
		pool: pool,
	}
}

// GetFriends gets the friends of a user.
func (s *FriendStore) GetFriends(userID uuid.UUID) ([]*types.Friend, error) {
	var friends []*types.Friend
	friends = []*types.Friend{}
	db, err := s.pool.Get()
	if err != nil {
		return friends, err
	}
	defer s.pool.Release()

	rows, err := db.Query(queries.GetFriends, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return friends, interfaces.ErrFriendNotFound
		}
		return friends, err
	}
	defer rows.Close()

	for rows.Next() {
		friend := &types.Friend{}
		err = rows.Scan(&friend.RID, &friend.UID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Bio, &friend.Avatar, &friend.AcceptedAt)
		if err != nil {
			return friends, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

// GetFriend gets a friend by its id.
func (s *FriendStore) GetFriend(userID uuid.UUID, friendID uuid.UUID) (*types.Friend, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	friend := &types.Friend{}
	err = db.QueryRow(queries.GetFriend, userID, friendID).Scan(&friend.RID, &friend.UID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Bio, &friend.Avatar, &friend.AcceptedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrFriendNotFound
		}
		return nil, err
	}
	return friend, nil
}

// DeleteFriend deletes a friend by its id.
func (s *FriendStore) DeleteFriend(userID, fID uuid.UUID) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	a, err := db.Exec(queries.DeleteFriend, userID, fID, userID, fID)
	if err != nil {
		return err
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return interfaces.ErrFriendNotFound
	}
	return nil
}

// GetFriendRequests gets the friend requests of a user.
func (s *FriendStore) GetFriendRequests(userID uuid.UUID) ([]*types.Friend, error) {
	var friends []*types.Friend
	db, err := s.pool.Get()
	if err != nil {
		return friends, err
	}
	defer s.pool.Release()

	rows, err := db.Query(queries.GetFriendRequests, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return friends, interfaces.ErrFriendNotFound
		}
		return friends, err
	}
	defer rows.Close()

	friends = []*types.Friend{}
	for rows.Next() {
		friend := &types.Friend{}
		err = rows.Scan(&friend.RID, &friend.UID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Bio, &friend.Avatar, &friend.AcceptedAt)
		if err != nil {
			return friends, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

// GetFriendRequest gets a friend request by its id.
func (s *FriendStore) GetFriendRequest(userID, uid uuid.UUID) (*types.Friend, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	friend := &types.Friend{}
	err = db.QueryRow(queries.GetFriendRequest, userID, uid).Scan(&friend.RID, &friend.UID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Bio, &friend.Avatar, &friend.AcceptedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrFriendNotFound
		}
		return nil, err
	}
	return friend, nil
}

// AcceptFriendRequest accepts a friend request.
func (s *FriendStore) AcceptFriendRequest(userUID, uid uuid.UUID) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	a, err := db.Exec(queries.AcceptFriendRequest, userUID, userUID, uid)
	if err != nil {
		return err
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return interfaces.ErrFriendNotFound
	}
	return nil
}

// DeleteFriendRequest deletes a friend request by its id.
func (s *FriendStore) DeleteFriendRequest(userID, reqID uuid.UUID) error {
	db, err := s.pool.Get()
	if err != nil {
		return err
	}
	defer s.pool.Release()

	a, err := db.Exec(queries.DeleteFriendRequest, userID, userID, reqID)
	if err != nil {
		return err
	}

	if n, err := a.RowsAffected(); err != nil || n == 0 {
		return interfaces.ErrFriendNotFound
	}

	return nil
}

// GetFriendsStatus gets the friends status of a user.
func (s *FriendStore) GetFriendsStatus(userID uuid.UUID) ([]*types.FriendStatus, error) {
	var friends []*types.FriendStatus
	friends = []*types.FriendStatus{}
	db, err := s.pool.Get()
	if err != nil {
		return friends, err
	}
	defer s.pool.Release()

	rows, err := db.Query(queries.GetFriendsStatus, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return friends, interfaces.ErrFriendNotFound
		}
		return friends, err
	}
	defer rows.Close()

	for rows.Next() {
		friend := &types.FriendStatus{}
		err = rows.Scan(&friend.RID, &friend.UID, &friend.StatusID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Avatar, &friend.Title, &friend.ResourceURI, &friend.ResourceThumbnail)
		if err != nil {
			return friends, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

// GetFriendStatus gets the friend status of a user.
func (s *FriendStore) GetFriendStatus(userID uuid.UUID, friendID uuid.UUID) (*types.FriendStatus, error) {
	db, err := s.pool.Get()
	if err != nil {
		return nil, err
	}
	defer s.pool.Release()

	friend := &types.FriendStatus{}
	err = db.QueryRow(queries.GetFriendStatus, userID, friendID).Scan(&friend.RID, &friend.UID, &friend.StatusID, &friend.Username, &friend.FirstName, &friend.LastName, &friend.Avatar, &friend.Title, &friend.ResourceURI, &friend.ResourceThumbnail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrFriendStatusNotFound
		}
		return nil, err
	}
	return friend, nil
}
