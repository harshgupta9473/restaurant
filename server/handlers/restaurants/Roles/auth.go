package roles

import (
	"encoding/json"
	"net/http"


	rolesHelper "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/roles"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

// RequestRole allows a user to request a role in a restaurant while being common user
func RequestRole(w http.ResponseWriter, r *http.Request) {
	var req roleModels.NewUserRequestForRole

	user,err:=middlewares.GetUserContext(r)
	if err!=nil{
		//
		return
	}
	userId:=user.UserID
	

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Error:   "INVALID_REQUEST_DATA",
		})
		return
	}

	_, err =rolesHelper.GetRoleByRoleIdandCode(req.RoleId,req.Code,req.RestaurantID)
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
			Status:  "error",
			Message: "Role not found",
			Error:   "ROLE_NOT_FOUND",
		})
		return
	}

	err = rolesHelper.InsertIntoStaffTable(req.RoleId,req.RestaurantID,userId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to request role",
			Error:   "ROLE_REQUEST_FAILED",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Role requested successfully",
	})
}

func ApproveRoleRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       int64 `json:"user_id"`
		RestaurantID int64 `json:"restaurant_id"`
		RoleID       int64 `json:"role_id"`
		Approve      bool  `json:"approve"`
	}

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Error:   "INVALID_REQUEST_DATA",
		})
		return
	}
	user,err:=middlewares.GetRestaurantUserContext(r)
	if err!=nil{
		//
		return
	}
	// 🔐 Extract acting user from context (set by JWT middleware)
	actingUserID := user.UserID
	actingRoleID := user.RoleId

	// ✅ 1. Check if acting user is Owner or has permission to approve roles for target role
	hasPermission, err := database.HasApprovePermission(actingRoleID, req.RoleID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Authorization check failed",
			Error:   "AUTH_CHECK_FAILED",
		})
		return
	}

	if !hasPermission {
		utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
			Status:  "error",
			Message: "You are not authorized to approve this role",
			Error:   "NOT_AUTHORIZED",
		})
		return
	}

	// ✅ 2. Proceed with approval/rejection
	err = database.UpdateRoleApproval(req.UserID, req.RestaurantID, req.RoleID, req.Approve)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to approve/reject role",
			Error:   "APPROVAL_FAILED",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Role request approved/rejected",
	})
}
