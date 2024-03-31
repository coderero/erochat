package queries

// SQL queries template constants for user.

const (
	// GetUser returns a user by uid.
	GetUserByID = `SELECT * FROM users WHERE id = ?`

	// GetUser returns a user by uid.
	GetUserByUID = `SELECT * FROM users WHERE uid = ?`

	// GetUser returns a user by username.
	GetUserByEmail = `SELECT * FROM users WHERE email = ?`

	// GetUser returns a user by username.
	GetUserByUsername = `SELECT * FROM users WHERE username = ?`

	// CreateUser creates a new user.
	CreateUser = `INSERT INTO users (uid,username, email, password) VALUES (UUID(),?, ?, ?)`

	// UpdateUser updates a user.
	UpdateUser = `UPDATE users SET username = COALESCE(?, username), email = COALESCE(?, email), password = COALESCE(?, password), updated_at = now() WHERE uid = ?`

	// DeleteUser deletes a user.
	DeleteUser = `UPDATE users SET deleted_at = now() WHERE uid = ?`
)
