package restaurants

import (
	"encoding/json"
	"net/http"
	"github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/Registration"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func RestaurantAccountRequest(w http.ResponseWriter, r *http.Request) {
	
	var req models.RestaurantFormReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Request",
			Error:   "Not the Correct Format",
		})
		return
	}

	
	restaurantID, err := restaurants.CreateRestaurantRequest(&req)

	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to create restaurant",
			Error:   err.Error(),
		})
		return
	}


	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Restaurant account created successfully",
		Data:    map[string]int64{"restaurant_id": restaurantID},
	})
}


