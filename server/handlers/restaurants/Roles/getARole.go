package roles

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	rolesHelper "github.com/harshgupta9473/restaurantmanagement/helpers/restaurants/roles"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	roleModels "github.com/harshgupta9473/restaurantmanagement/models/roles"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

const (
	IsApprovedPending  = 0
	IsApprovedAccepted = 1
	IsApprovedRejected = -1
	IsApprovedBlocked  = -2
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
	handleStaffStatusChange(w, r, IsApprovedAccepted, "Role request approved successfully", "Approval Request Failed")

}

// reject request   is_approved=-1
func RejectRoleRequest(w http.ResponseWriter, r *http.Request) {
	handleStaffStatusChange(w, r, IsApprovedRejected, "Staff request rejected", "Not able to Reject the Staff")

}

// block request  is_approved=-2
func BlockStaffRequest(w http.ResponseWriter, r *http.Request) {
	handleStaffStatusChange(w, r, IsApprovedBlocked, "Staff blocked successfully", "Blocking Failed")
}

func handleStaffStatusChange(w http.ResponseWriter, r *http.Request, newStatus int, successMsg, failMsg string) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	vars := mux.Vars(r)
	roleID, err := strconv.ParseInt(vars["roleID"], 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid roleID",
			Error:   "INVALID_ROLE_ID",
		})
		return
	}

	staffId, err := strconv.ParseInt(vars["staffID"], 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid staffID",
			Error:   "INVALID_STAFF_ID",
		})
		return
	}

	result, err := rolesHelper.ChangeStatusOfIsApprovedInStaffTable(newStatus, user.UserID, staffId, user.RestaurantId, roleID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: failMsg,
			Error:   "STATUS_UPDATE_FAILED",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "info",
			Message: "No changes made. Already updated or invalid request.",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: successMsg,
	})
}

func GetAllRolesThatsUnderAuthority(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	restaurantId := user.RestaurantId
	userID:=user.UserID

	result,err:=rolesHelper.GetAllRolesThatsUnderAuthorityHelper(userID,restaurantId)

if  err != nil {
	utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
		Status:  "error",
		Message: err.Error(),
		Error:   "DB_ERROR",
	})
	return
}

utils.WriteJson(w, http.StatusOK, utils.APIResponse{
	Status:  "success",
	Message: "Roles fetched successfully",
	Data:    result,
})

}

func GetAllRequestInRoleThatsUnderAuthority(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetRestaurantUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized access",
			Error:   "USER_CONTEXT_MISSING",
		})
		return
	}

	restaurantId := user.RestaurantId
	vars := mux.Vars(r)

	roleID, err := strconv.ParseInt(vars["roleID"], 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid roleID",
			Error:   "INVALID_ROLE_ID",
		})
		return
	}

	result,err:=rolesHelper.GetAllRequestInRoleThatsUnderAuthorityHelper(roleID,restaurantId)

	if err!=nil{
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Error:   "DB_ERROR",
		})
		return
	}

	if len(result) == 0 {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "No pending staff requests found",
			Data:    []roleModels.PersonWithRoles{},
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Pending staff requests fetched successfully",
		Data:    result,
	})
}
