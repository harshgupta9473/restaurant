package roles

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/harshgupta9473/restaurantmanagement/db"
// 	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
// 	"github.com/harshgupta9473/restaurantmanagement/utils"
// )

// func CreateRoles(w http.ResponseWriter, r *http.Request) {
// 	var rolereq roleModels.RoleReq
// 	err := json.NewDecoder(r.Body).Decode(&rolereq)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}


// 	conn := db.GetDB()

// 	query := `
// 		INSERT INTO roles (name, description, restaurant_id, level, is_global, is_custom)
// 		VALUES ($1, $2, $3, $4, FALSE, TRUE)
// 		RETURNING id;
// 	`

// 	var result roleModels.Roles
// 	err = conn.QueryRow(query, rolereq.Name,rolereq.Description, rolereq.RestaurantId, rolereq.Level).Scan(&result.ID)
// 	if err != nil {
// 		log.Println("Error inserting role:", err)
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Could not create role", http.StatusInternalServerError)
// 		} else {
// 			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
// 		}
// 		return
// 	}
// 	result.Description=re
// 	utils.WriteJson(w,http.StatusCreated,utils.APIResponse{
// 		Status: "success",
// 		Message: "Created the Role",
// 		Data:  map[string]int64{"role_id":roleId},
// 	})
// }




