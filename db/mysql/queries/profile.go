package queries

// SQL queries template constants for user profile.
const (
	// GetUserProfile returns a profile by uid.
	GetUserProfileByID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.bio, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.id = ?`

	// GetUserProfile returns a profile by uid.
	GetUserProfileByUID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.bio, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.uid = ?`

	// GetUserProfile returns a profile by user_id.
	GetUserProfileByUserID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.bio, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.user_id = ?`

	// GetUserProfile returns a profile by email.
	GetUserProfileByEmail = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.bio, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE u.email = ?`

	// CreateUserProfile creates a new profile.
	CreateUserProfile = `INSERT INTO profiles (uid, user_id, first_name, last_name, bio, avatar) VALUES (?, ?, ?, ?, ?, ?)`

	// UpdateUserProfile updates a profile.
	CreateFriendship = `INSERT INTO friendships (uid, user1, user2) VALUES (UUID(), ?, ?)`

	// UpdateUserProfile updates a profile.
	UpdateUserProfile = `UPDATE profiles SET first_name = COALESCE(?, first_name), last_name = COALESCE(?, last_name), bio = COALESCE(?, bio), avatar = COALESCE(?, avatar), updated_at = now() WHERE user_id = ?`

	// DeleteUserProfile deletes a profile.
	DeleteUserProfile = `UPDATE profiles SET deleted_at = now() WHERE user_id = ?`

	// ReactivateUserProfile reactivates a profile.
	ReactivateUserProfile = `UPDATE profiles SET deleted_at = NULL WHERE user_id = ?`
)
