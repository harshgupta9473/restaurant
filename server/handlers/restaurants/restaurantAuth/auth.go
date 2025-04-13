package restaurants

import (
	"database/sql"
	"encoding/json"
	"net/http"

	helpers "github.com/harshgupta9473/restaurantmanagement/helpers/auth"
	middlewaresHelper "github.com/harshgupta9473/restaurantmanagement/helpers/middleware"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"

	commonModels "github.com/harshgupta9473/restaurantmanagement/models/common"
	models "github.com/harshgupta9473/restaurantmanagement/models/user"
	"github.com/harshgupta9473/restaurantmanagement/utils"
	"golang.org/x/crypto/bcrypt"
)

func LoginHelper(w http.ResponseWriter, r *http.Request, email, password, role string) (*models.User, bool) {
	user, err := helpers.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: "Invalid Request",
				Error:   "Invalid Email or Password",
			})
		} else {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   "DB_ERROR",
			})
		}
		return nil, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Email Or Password",
			Error:   "Invalid_EMAIL_PASSWORD",
		})
		return nil, false
	}
	return user, true
}
func RestaurantOwnerLogIN(w http.ResponseWriter, r *http.Request) {
	var req commonModels.RoleLoginReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	
	if req.Role != "owner" {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid role provided",
			Error:   "Access restricted to restaurant owners only",
		})
		return
	}

	
	user, yes := LoginHelper(w, r, req.Email, req.Password, req.Role)
	if !yes {
		return 
	}


	exists, err := middlewaresHelper.IsRestaurantOwner(user.Id)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not verify restaurant ownership",
			Error:   err.Error(),
		})
		return
	}
	if !exists {
		utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
			Status:  "error",
			Message: "You are not registered as a restaurant owner",
			Error:   "Unauthorized",
		})
		return
	}

	
	accessToken, refreshToken, err := middlewares.GenerateTokenForRole(user.Id, user.Verified, req.Role)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not generate token",
			Error:   err.Error(),
		})
		return
	}

	
	utils.SetCookie(w, "access_token", accessToken, 24*3600)
	utils.SetCookie(w, "refresh_token", refreshToken, 7*24*3600)

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Login successful",
	})
}

