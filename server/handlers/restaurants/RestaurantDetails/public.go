package restaurants

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	restaurants "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/RestaurantDetails"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func AboutRestaurant(w http.ResponseWriter, r *http.Request) {

	restaurantIdStr := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdStr, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
			Error:   err.Error(),
		})
		return
	}

	details, err := restaurants.GetCompletePublicRestaurantDetails(restaurantId)
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
		Message: "Restaurant details fetched successfully",
		Data:    details,
	})
}

// e.g.  /api/reviews/12?page=2&limit=5
func GetReview(w http.ResponseWriter, r *http.Request) {
	restaurantIdStr := mux.Vars(r)["id"]
	restaurantId, err := strconv.ParseInt(restaurantIdStr, 10, 64)

	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid restaurant ID",
		})
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	reviews, err := restaurants.GetReview(restaurantId, page, limit)
	if err != nil && err != sql.ErrNoRows {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch reviews",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Fetched reviews",
		Data: map[string]interface{}{
			"reviews": reviews,
			"page":    page,
			"limit":   limit,
		},
	})
}


