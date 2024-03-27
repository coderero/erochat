package interfaces

import (
	"github.com/coderero/erochat-server/types"
	"github.com/google/uuid"
)

type FriendStore interface {
	// GetFriends gets the friends of a user.
	GetFriends(userID uuid.UUID) ([]*types.User, error)

	// CreateFriend creates a new friend.
	GetFriend(userID uuid.UUID, friendID uuid.UUID) (*types.User, error)

	// DeleteFriend deletes a friend by its id.
	DeleteFriend(userID uuid.UUID, friendID uuid.UUID) error

	// GetFriendRequests gets the friend requests of a user.
	GetFriendRequests(userID uuid.UUID) ([]*types.User, error)

	// GetFriendRequest gets a friend request by its id.
	GetFriendRequest(userID uuid.UUID, friendID uuid.UUID) (*types.User, error)

	// DeleteFriendRequest deletes a friend request by its id.
	DeleteFriendRequest(userID uuid.UUID, friendID uuid.UUID) error

	// GetFriendsStatus gets the friends status of a user.
	GetFriendsStatus(userID uuid.UUID) ([]*types.UserStatus, error)

	// GetFriendStatus gets the friend status of a user.
	GetFriendStatus(userID uuid.UUID, friendID uuid.UUID) (*types.UserStatus, error)
}
