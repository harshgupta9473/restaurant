package restaurants

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	restaurants "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/RestaurantDetails"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func GetRestaurantsPrivateDetails(w http.ResponseWriter, r *http.Request) {
	// Get restaurant ID from path variable
	restaurantIdString := mux.Vars(r)["restaurant_id"]
	restaurantId, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}

	details, err := restaurants.GetDetailsOfRestaurant(restaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch restaurant details",
			Error:   err.Error(),
		})
		return
	}

	
	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Fetched restaurant private details successfully",
		Data:    details,
	})
}

