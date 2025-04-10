package restaurants

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/restaurantmanagement/helpers/restaurants"
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


func ApproveRestaurantRequest(w http.ResponseWriter, r *http.Request) {
	restaurantIdString := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}

	restaurant, err := restaurants.ApproveAndCreateRestaurant(restaurantId)
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
		Data:    map[string]int64{"restaurant_id": restaurant.ID},
	})
}

func BlockRestaurantRequest(w http.ResponseWriter, r *http.Request) {
	restaurantIdString := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}

	res, err := restaurants.BlockRestaurantRequest(restaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to block request",
			Error:   err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "No pending request found to block",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Restaurant request blocked successfully",
	})
}


func DeleteRestaurantRequest(w http.ResponseWriter, r *http.Request) {
	restaurantIdString := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}


	res, err := restaurants.DeleteBlockedRestaurantRequest(restaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to delete request",
			Error:   err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "No pending request found to delete",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Restaurant request deleted successfully",
	})
}


func ReRequestRestaurantRequest(w http.ResponseWriter, r *http.Request) {
	restaurantIdString := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}


	res, err := restaurants.DeletePendingRestaurantRequest(restaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to delete blocked/rejected request",
			Error:   err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "No blocked or rejected request found to delete",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Old request deleted, you can now re-request",
	})
}

