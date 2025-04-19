package middlewaresHelper

import (

	"github.com/harshgupta9473/restaurantmanagement/db"
)

// AuthorityPermissionCheck verifies whether a user (actingUserID) has a specific permission
// over a target role (targetRoleID) within a restaurant (restaurantID).
// It does this by checking:
// - if the user has an approved staff role in the given restaurant,
// - if that role is mapped to the specified permission name,
// - and if that permission allows acting on the target role.
// Returns true if permission exists, false with an error otherwise.

func AuthorityPermissionCheck(userId, restaurantId int64, targetRoleId int64, permission string) (error) {
	conn := db.GetDB()

	query := `
	SELECT 1
	FROM staff_roles AS sr
	JOIN role_target_permissions AS rtp ON sr.role_id = rtp.actor_role_id
	JOIN permissions AS p ON rtp.permission_id = p.id
	WHERE sr.user_id = $1
	  AND sr.restaurant_id = $2
	  AND sr.is_approved = TRUE
	  AND p.name = $3
	  AND rtp.target_role_id = $4
	LIMIT 1;
	`

	var exists int
	err := conn.QueryRow(query, userId, restaurantId, permission, targetRoleId).Scan(&exists)
	if err != nil {
		 return err
	}

	return nil
}




// CheckPermissionFunc checks if a user has a general permission in a restaurant,
// without targeting a specific role.

func CheckPermissionFunc(userID, restaurantID int64, permission string)(error){

	conn:=db.GetDB()

	query:=`
	SELECT 1 
	FROM staff_roles as sr
	JOIN role_permissions as rp ON rp.role_id=sr.role_id
	JOIN permissions as p ON p.id=rp.permission_id
	WHERE sr.user_id=$1
	AND sr.is_approved=1
	AND sr.restaurant_id=$2
	AND p.name=$3
	LIMIT 1;
	`
	var exists int
	err:=conn.QueryRow(query,userID,restaurantID,permission).Scan(&exists)

	if err!=nil{
		return err
	}

	return nil
}