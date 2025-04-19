package middlewaresHelper

import (
	"fmt"

	"github.com/harshgupta9473/restaurantmanagement/db"
)

func IsVerified(user_id int64) error {
	db := db.GetDB()
	var verified bool
	err := db.QueryRow(`SELECT verified FROM users WHERE id=$1`, user_id).Scan(&verified)
	if err != nil {
		return fmt.Errorf("Failed to fetch user verification status")
	}

	if !verified {
		return fmt.Errorf("Please verify your account to proceed")
		// fmt.Errorf()
	}
	return nil
}
