package roles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/harshgupta9473/restaurantmanagement/db"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func CreateRoles(w http.ResponseWriter, r *http.Request) {
	var role roleModels.RoleReq
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}


	conn := db.GetDB()

	query := `
		INSERT INTO roles (name, description, restaurant_id, level, is_global, is_custom)
		VALUES ($1, $2, $3, $4, FALSE, TRUE)
		RETURNING id;
	`

	var roleId int64
	err = conn.QueryRow(query, role.Name, role.Description, role.RestaurantId, role.Level).Scan(&roleId)
	if err != nil {
		log.Println("Error inserting role:", err)
		if err == sql.ErrNoRows {
			http.Error(w, "Could not create role", http.StatusInternalServerError)
		} else {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		}
		return
	}
	utils.WriteJson(w,http.StatusCreated,utils.APIResponse{
		Status: "success",
		Message: "Created the Role",
		Data:  {

		}
	})
}
