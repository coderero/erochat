package interfaces

import (
	"errors"

	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

var (
	// ErrFriendNotFound is an error that is returned when the friend is not found.
	ErrFriendNotFound = errors.New("friend not found")

	// ErrFriendStatusNotFound is an error that is returned when the friend status is not found.
	ErrFriendStatusNotFound = errors.New("friend status not found")
)

type FriendStore interface {
	// GetFriends gets the friends of a user.
	GetFriends(userID uuid.UUID) ([]*types.Friend, error)

	// CreateFriend creates a new friend.
	GetFriend(userID uuid.UUID, friendID uuid.UUID) (*types.Friend, error)

	// DeleteFriend deletes a friend by its id.
	DeleteFriend(userID, fID uuid.UUID) error
	// GetFriendRequests gets the friend requests of a user.
	GetFriendRequests(userID uuid.UUID) ([]*types.Friend, error)

	// GetFriendRequest gets a friend request by its id.
	GetFriendRequest(userID, uid uuid.UUID) (*types.Friend, error)

	// AcceptFriendRequest accepts a friend request.
	AcceptFriendRequest(userUID, uid uuid.UUID) error

	// DeleteFriendRequest deletes a friend request by its id.
	DeleteFriendRequest(userID, reqID uuid.UUID) error

	// GetFriendsStatus gets the friends status of a user.
	GetFriendsStatus(userID uuid.UUID) ([]*types.FriendStatus, error)

	// GetFriendStatus gets the friend status of a user.
	GetFriendStatus(userID uuid.UUID, friendID uuid.UUID) (*types.FriendStatus, error)
}
