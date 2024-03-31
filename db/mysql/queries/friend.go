package queries

// SQL queries template constants for friend.
const (
	// GetFriends returns all friends of a user.
	GetFriends = `CALL get_friends_or_requests(?, false)`

	// GetFriend returns a friend by its id.
	GetFriend = `CALL get_friend(?, ?)`

	// DeleteFriend deletes a friend by its id.
	DeleteFriend = `DELETE FROM friendships WHERE (user1 = ? OR user2 = ?) AND (user2 = ? OR user1 = ?) AND accepted = true`

	// GetFriendRequests returns all friend requests of a user.
	GetFriendRequests = `CALL get_friends_or_requests(?, true)`

	// GetFriendRequest returns a friend request by its id.
	GetFriendRequest = `CALL get_friend_request(?, ?)`

	// AcceptFriendRequest accepts a friend request.
	AcceptFriendRequest = `UPDATE friendships SET accepted = true, accepted_at = now() WHERE (user1 = ? OR user2 = ?) AND uid = ? AND accepted = false`

	// DeleteFriendRequest deletes a friend request by its id.
	DeleteFriendRequest = `DELETE FROM friendships WHERE (user1 = ? OR user2 = ?) AND (uid = ?) AND (accepted = false)`

	// GetFriendsStatus returns all friends status of a user.
	GetFriendsStatus = `CALL get_friends_status(?)`

	// GetFriendStatus returns a friend status of a user.
	GetFriendStatus = `CALL get_friend_status(?, ?)`
)
