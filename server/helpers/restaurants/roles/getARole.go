package rolesHelper

import (
	"database/sql"
	"fmt"

	"github.com/harshgupta9473/restaurantmanagement/db"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
)

func GetRoleByRoleIdandCode(roleId int64, code string, restaurantId int64) (roleModels.Roles, error) {
	var role roleModels.Roles
	conn := db.GetDB()
	query := `
		SELECT id, code, name, description, restaurant_id, is_global, level, is_custom, created_at, updated_at 
		FROM roles
		WHERE id = $1 AND code = $2 AND restaurant_id=$3
	`

	row := conn.QueryRow(query, roleId, code, restaurantId)

	err := row.Scan(&role.ID, &role.Code, &role.Name, &role.Description, &role.RestaurantId,
		&role.Is_Global, &role.Level, &role.Is_Custom, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {

		if err == sql.ErrNoRows {
			return role, fmt.Errorf("role not found with ID %d and code %s", roleId, code)
		}
		return role, fmt.Errorf("error querying database: %v", err)
	}
	return role, nil
}

func ChangeStatusOfIsApprovedInStaffTable(is_approved int, actingUserID int64, staffId int64, restaurantId int64, roleID int64) (sql.Result, error) {
	query := `
		UPDATE staff_table
		SET is_approved = $1,
		    manager_id = $2
		WHERE id = $3 AND restaurant_id = $4 AND role_id = $5
	`

	conn := db.GetDB()
	result, err := conn.Exec(query, is_approved, actingUserID, staffId, restaurantId, roleID)
	return result, err
}

func InsertIntoStaffTable(roleID, restaurantId, userId int64) error {
	conn := db.GetDB()
	query := `INSERT INTO staff_table(restaurant_id,user_id,role_id,is_approved)
	VALUES ($1,$2,$3,$4)
	`
	_, err := conn.Exec(query, restaurantId, userId, roleID)
	return err
}
