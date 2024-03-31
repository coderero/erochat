package queries

// SQL queries template constants for user status.
const (
	// GetStatus returns a status by uid which is created under 24 hours.
	GetStatus = `SELECT id, uid, user_uid, title, resource_uri, resource_thumbnail, created_at FROM status WHERE uid = ? AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR) AND deleted_at IS NULL`

	// GetStatusByID returns a status by id.
	GetStatusByID = `SELECT id, uid, user_uid, title, resource_uri, resource_thumbnail, created_at FROM status WHERE (user_uid = ?, uid = ?) AND deleted_at IS NULL`

	// GetUsersStatus returns all status of a user.
	GetUsersStatus = `SELECT id, uid, user_uid, title, resource_uri, resource_thumbnail, created_at FROM status WHERE user_uid = ? AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR) AND deleted_at IS NULL`

	// CreateStatus creates a new status.
	CreateStatus = `INSERT INTO status (uid, user_uid, title, resource_uri, resource_thumbnail) VALUES (UUID(),?,?,?,?)`

	// DeleteStatus deletes a status by uid.
	DeleteStatus = `UPDATE status SET deleted_at = now() WHERE (user_uid = ? AND uid = ?) AND deleted_at IS NULL AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR)`
)
