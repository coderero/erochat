package queries

// SQL queries template constants for user profile.
const (
	GetUserProfileByID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.id = ?`

	GetUserProfileByUID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.uid = ?`

	GetUserProfileByUserID = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.user_id = ?`

	GetUserProfileByEmail = `SELECT p.id, p.uid, p.user_id, p.first_name, p.last_name, p.avatar, u.username, u.email, p.created_at, p.updated_at, p.deleted_at FROM profiles p JOIN users u ON p.user_id = u.id WHERE u.email = ?`

	CreateUserProfile = `INSERT INTO profiles (uid, user_id, first_name, last_name, avatar) VALUES (UUID(), ?, ?, ?, ?)`

	UpdateUserProfile = `UPDATE profiles SET first_name = COALESCE(?, first_name), last_name = COALESCE(?, last_name), avatar = COALESCE(?, avatar), updated_at = now() WHERE user_id = ?`

	DeleteUserProfile = `UPDATE profiles SET deleted_at = now() WHERE user_id = ?`
)
