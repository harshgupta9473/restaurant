package middlewares

import (
	"fmt"

	"github.com/harshgupta9473/restaurantmanagement/db"
)

func IsSuperAdmin(userID int64) (bool, error) {
	conn := db.GetDB()

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM super_admins WHERE user_id = $1)`
	err := conn.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check super admin status: %v", err)
	}
	return exists, nil
}
