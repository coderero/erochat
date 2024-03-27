package queries

// SQL queries template constants for user status.
const (
	// GetStatus returns a status by uid which is created under 24 hours.
	GetStatus = `SELECT id, uid, user_id, title, resource_uri, resource_thumbnail, created_at FROM statuses WHERE uid = ? AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR) AND deleted_at IS NULL`

	// GetStatusByID returns a status by id.
	GetStatusByID = `SELECT id, uid, user_id, title, resource_uri, resource_thumbnail, created_at FROM statuses WHERE id = ? AND deleted_at IS NULL`

	// GetUsersStatus returns all statuses of a user.
	GetUsersStatus = `SELECT id, uid, user_id, title, resource_uri, resource_thumbnail, created_at FROM statuses WHERE user_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR) AND deleted_at IS NULL`

	// CreateStatus creates a new status.
	CreateStatus = `INSERT INTO statuses (uid, user_id, title, resource_uri, resource_thumbnail) VALUES (UUID(),?,?,?,?)`

	// DeleteStatus deletes a status by uid.
	DeleteStatus = `UPDATE statuses SET deleted_at = now() WHERE (uid = ? AND user_id = ?) AND deleted_at IS NULL AND created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR)`
)
