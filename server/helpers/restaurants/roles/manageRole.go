package rolesHelper

import (
	"github.com/harshgupta9473/restaurantmanagement/db"
)

func InsertIntoRoles(name string,description string,restaurantId int64,level int,mangerId *int64) (int64, error) {
	conn := db.GetDB()
	query := `
		INSERT INTO roles (name, description, restaurant_id, level, is_global, is_custom,manager_id)
		VALUES ($1, $2, $3, $4, FALSE, TRUE)
		RETURNING id;
	`
	var id int64

	err := conn.QueryRow(query,name, description, restaurantId,level,mangerId).Scan(id)
	if err!=nil{
		return 0,err
	}
	return id,nil
}
