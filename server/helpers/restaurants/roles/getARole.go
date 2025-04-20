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
	VALUES ($1,$2,$3,$4,$5)
	`
	_, err := conn.Exec(query, restaurantId, userId, roleID)
	return err
}

func GetAllRolesThatsUnderAuthorityHelper(userID, restaurantId int64, permissionId int64) ([]roleModels.Roles, error) {

	query := `
SELECT r.id, r.code, r.name, r.description, r.restaurant_id, r.is_global,r.level, 
       r.is_custom, r.manager_id, r.created_at, r.updated_at
FROM staff_roles as sr 
JOIN role_target_permissions as rtp ON rtp.actor_role_id = sr.role_id
JOIN roles as r ON r.id = rtp.target_role_id
WHERE sr.user_id = $1
  AND sr.restaurant_id = $2
  AND sr.is_approved = 1
  AND rtp.action_permission_id=$3
`

	rows, err := db.GetDB().Query(query, userID, restaurantId, permissionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []roleModels.Roles

	for rows.Next() {
		var role roleModels.Roles
		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.RestaurantId,
			&role.Is_Global,
			&role.Level,
			&role.Is_Custom,
			&role.ManagerId,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return result, err
		}
		result = append(result, role)
	}

	if err := rows.Err(); err != nil {

		return nil, err
	}
	return result, nil
}

func GetAllStaffMemberOfRoleThatsUnderAuthorityHelper(
	userID, roleID, restaurantId, permissionID int64,
) ([]roleModels.PersonWithRoles, error) {

	query := `
	SELECT sr.id, sr.user_id, sr.restaurant_id, sr.role_id, r.name, 
	       u.first_name, u.middle_name, u.last_name, u.manager_id, u.blocked, 
	       sr.created_at, sr.updated_at, sr.is_approved
	FROM users AS u
	JOIN staff_roles AS sr ON u.id = sr.user_id
	JOIN roles AS r ON sr.role_id = r.id
	JOIN role_target_permissions AS rtp ON rtp.target_role_id = sr.role_id
	WHERE sr.role_id = $1
	  AND sr.restaurant_id = $2
	  AND rtp.action_permission_id = $3
	  AND EXISTS (
	      SELECT 1
	      FROM staff_roles AS actor_sr
	      WHERE actor_sr.user_id = $4
	        AND actor_sr.restaurant_id = $2
	        AND actor_sr.is_approved = TRUE
	        AND actor_sr.role_id = rtp.actor_role_id
	  )
	`

	rows, err := db.GetDB().Query(query, roleID, restaurantId, permissionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var result []roleModels.PersonWithRoles

	for rows.Next() {
		var person roleModels.PersonWithRoles
		err := rows.Scan(
			&person.StaffId,
			&person.UserID,
			&person.RestaurantId,
			&person.RoleId,
			&person.RoleName,
			&person.FirstName,
			&person.MiddleName,
			&person.LastName,
			&person.ManagerId,
			&person.Blocked,
			&person.CreatedAt,
			&person.UpdatedAt,
			&person.Is_Approved,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}
