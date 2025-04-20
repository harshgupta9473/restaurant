package roles

import (
	"encoding/json"
	"net/http"

	"github.com/harshgupta9473/restaurantmanagement/db"
	helpers "github.com/harshgupta9473/restaurantmanagement/helpers/auth"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
	"golang.org/x/crypto/bcrypt"
)

func RoleLogin(w http.ResponseWriter, r *http.Request) {
	var req roleModels.RestaurantRoleLogin

	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}


	user, err := helpers.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}
	if !user.Verified{
		//
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}


	query := `SELECT 1 FROM staff_roles WHERE user_id = $1 AND restaurant_id = $2`
	var exists int
	if err := db.GetDB().QueryRow(query, user.Id, req.RestaurantId).Scan(&exists); err != nil {
		utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
			Status:  "error",
			Message: "You are not authorized for this restaurant",
		})
		return
	}

	restaurantAccessToken, err := middlewares.GenerateRestaurantAccessToken(user.Id, req.RestaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to generate access token",
		})
		return
	}

	restaurantRefreshToken, err := middlewares.GenerateRestaurantRefreshToken(user.Id, req.RestaurantId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	
	utils.SetCookie(w, "restaurant_access_token", restaurantAccessToken, 24*3600)
	utils.SetCookie(w, "restaurant_refresh_token", restaurantRefreshToken, 7*24*3600)

	
	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Login successful",
	})
}
