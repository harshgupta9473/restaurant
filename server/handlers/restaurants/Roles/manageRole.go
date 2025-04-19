package roles

import (
	"encoding/json"

	"log"
	"net/http"
	"time"

	rolesHelper "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/roles"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

// for owner
func CreateRolesByOwner(w http.ResponseWriter, r *http.Request) {
	var rolereq roleModels.NewRoleCreationRequestOwner
	if err := json.NewDecoder(r.Body).Decode(&rolereq); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   "INVALID_BODY",
		})
		return
	}

	id, err := rolesHelper.InsertIntoRoles(rolereq.Name, rolereq.Description, rolereq.RestaurantId, int(rolereq.Level))
	if err != nil {
		log.Println("Error inserting role:", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not create role",
			Error:   err.Error(),
		})
		return
	}

	var result roleModels.Roles
	result.ID = id
	result.Name = rolereq.Name
	result.Description = rolereq.Description
	result.RestaurantId = rolereq.RestaurantId
	result.Level = rolereq.Level
	result.Is_Global = false
	result.Is_Custom = true
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Created the role successfully",
		Data:    result,
	})
}



// for  user who has power to create roles
func CreateRolesWhoHavePowerToCreate(w http.ResponseWriter, r *http.Request) {
	var rolereq roleModels.NewRoleCreationRequestPower
	if err := json.NewDecoder(r.Body).Decode(&rolereq); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   "INVALID_BODY",
		})
		return
	}
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {

	}
	restaurantId := user.RestaurantId
	id, err := rolesHelper.InsertIntoRoles(rolereq.Name, rolereq.Description, restaurantId, int(rolereq.Level))
	if err != nil {
		log.Println("Error inserting role:", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not create role",
			Error:   err.Error(),
		})
		return
	}

	var result roleModels.Roles
	result.ID = id
	result.Name = rolereq.Name
	result.Description = rolereq.Description
	result.RestaurantId = restaurantId
	result.Level = rolereq.Level
	result.Is_Global = false
	result.Is_Custom = true
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Created the role successfully",
		Data:    result,
	})
}
