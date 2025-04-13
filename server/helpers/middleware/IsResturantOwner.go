package middlewares

import (
	"fmt"

	"github.com/harshgupta9473/restaurantmanagement/db"
)

func IsRestaurantOwner(user_id int64) (bool,error) {
	conn := db.GetDB()
	query := `Select EXISTS (Select 1 from restaurants where id=$1)`
	var IsOwner bool
	err := conn.QueryRow(query,user_id).Scan(&IsOwner)
	if err!=nil{
		return false,fmt.Errorf("failed to check stauts if owner %v",err)
	}
	return IsOwner,nil
}

