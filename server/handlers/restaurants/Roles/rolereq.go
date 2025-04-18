package roles

import (
	"encoding/json"
	"net/http"
	"strconv"

	rolesHelper "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/roles"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

// RequestRole allows a user to request a role in a restaurant while being common user   is_approved=0
func RequestRole(w http.ResponseWriter, r *http.Request) {
	var req roleModels.NewUserRequestForRole

	user, err := middlewares.GetUserContext(r)
	if err != nil {
		//
		return
	}
	userId := user.UserID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request data",
			Error:   "INVALID_REQUEST_DATA",
		})
		return
	}

	_, err = rolesHelper.GetRoleByRoleIdandCode(req.RoleId, req.Code, req.RestaurantID)
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
			Status:  "error",
			Message: "Role not found",
			Error:   "ROLE_NOT_FOUND",
		})
		return
	}

	err = rolesHelper.InsertIntoStaffTable(req.RoleId, req.RestaurantID, userId)
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

// approve request   is_approved=1
func ApproveRoleRequest(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	actingUserID := user.UserID
	restaurantId := user.RestaurantId

	
	roleID, err := strconv.ParseInt(r.URL.Query().Get("roleID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid role ID",
			Error:   "INVALID_ROLE_ID",
		})
		return
	}

	
	staffId, err := strconv.ParseInt(r.URL.Query().Get("staffID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid staff ID",
			Error:   "INVALID_STAFF_ID",
		})
		return
	}


	result, err:=rolesHelper.ChangeStatusOfIsApprovedInStaffTable(1,actingUserID,staffId,restaurantId,roleID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to approve role request",
			Error:   "APPROVAL_FAILED",
		})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "info",
			Message: "No changes made. Staff might already be approved or invalid request.",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Role request approved successfully",
	})
}


// reject request   is_approved=-1
func RejectRoleRequest(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	actingUserID := user.UserID
	restaurantId := user.RestaurantId

	roleID, err := strconv.ParseInt(r.URL.Query().Get("roleID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid roleID",
		})
		return
	}
	staffId, err := strconv.ParseInt(r.URL.Query().Get("staffID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid staffID",
		})
		return
	}

	
	result, err:=rolesHelper.ChangeStatusOfIsApprovedInStaffTable(-1,actingUserID,staffId,restaurantId,roleID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Rejection failed",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "info",
			Message: "No changes made. Already rejected or invalid request.",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Staff request rejected.",
	})
}


// block request  is_approved=-2
func BlockStaffRequest(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	actingUserID := user.UserID
	restaurantId := user.RestaurantId

	roleID, err := strconv.ParseInt(r.URL.Query().Get("roleID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid roleID",
		})
		return
	}

	staffId, err := strconv.ParseInt(r.URL.Query().Get("staffID"), 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid staffID",
		})
		return
	}

	

		result, err:=rolesHelper.ChangeStatusOfIsApprovedInStaffTable(-2,actingUserID,staffId,restaurantId,roleID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Blocking failed",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "info",
			Message: "Already blocked or invalid staff.",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Staff blocked successfully.",
	})
}


func GetAllRolesThatsUnderAuthority(w http.ResponseWriter,r *http.Request){

}

// by  role id
func GetAllRequestInRoleThatsUnderAuthority(w http.ResponseWriter,r *http.Request){

}

