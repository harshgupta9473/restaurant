package middlewares

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/harshgupta9473/restaurantmanagement/utils"
)


// General Permissions like "view_inventory etc"
func RequirePermission(permission string, checkPermissionFunc func(userID, restaurantID int64, permission string) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			

			user, err := GetRestaurantUserContext(r)
			if err != nil {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "Unauthorized: user context not found",
					Error:   "USER_CONTEXT_MISSING",
				})
				return
			}

			
			err = checkPermissionFunc(user.UserID, user.RestaurantId, permission)
			if err != nil {
				
				if err == sql.ErrNoRows {
					utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
						Status:  "error",
						Message: "Forbidden: insufficient permissions",
						Error:   "NO_PERMISSION",
					})
					return
				}

				// DB error
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Internal Server Error",
					Error:   "DB_ERROR",
				})
				return
			}

			
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuthority checks if the user has the required authority on a specific role.
func RequireAuthority(permission string, AuthorityPermissionCheck func(userId, restaurantId int64, targetRoleId int64, permission string) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			

			user, err := GetRestaurantUserContext(r)
			if err != nil {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "Unauthorized: user context not found",
					Error:   "USER_CONTEXT_MISSING",
				})
				return
			}

			
			roleID, err := strconv.ParseInt(r.URL.Query().Get("roleID"), 10, 64)
			if err != nil {
				utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
					Status:  "error",
					Message: "Bad Request: invalid roleID",
					Error:   "INVALID_ROLE_ID",
				})
				return
			}

			
			if err := AuthorityPermissionCheck(user.UserID, user.RestaurantId, roleID, permission); err != nil {
				
				if err == sql.ErrNoRows {
					utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
						Status:  "error",
						Message: "Forbidden: insufficient permissions",
						Error:   "NO_PERMISSION",
					})
					return
				}

				
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Internal Server Error",
					Error:   "DB_ERROR",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}