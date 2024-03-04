package queries

// SQL queries template constants for user.

const (
	GetUserByID = `SELECT * FROM users WHERE id = ?`

	GetUserByUID = `SELECT * FROM users WHERE uid = ?`

	GetUserByEmail = `SELECT * FROM users WHERE email = ?`

	GetUserByUsername = `SELECT * FROM users WHERE username = ?`

	CreateUser = `INSERT INTO users (uid,username, email, password) VALUES (UUID(),?, ?, ?)`

	UpdateUser = `UPDATE users SET username = COALESCE(?, username), email = COALESCE(?, email), password = COALESCE(?, password), updated_at = now() WHERE uid = ?`

	DeleteUser = `UPDATE users SET deleted_at = now() WHERE uid = ?`
)
