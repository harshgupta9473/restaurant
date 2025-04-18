package middlewares

import (
	"net/http"
	"strconv"

	"github.com/harshgupta9473/restaurantmanagement/utils"
)

// General Permissions like "view_inventory etc"
func RequirePermission(permission string, checkPermissionFunc func(userID, restaurantID int64, permission string) bool) func(http.Handler) http.Handler {
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

			if !checkPermissionFunc(user.UserID, user.RestaurantId, permission) {
				utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
					Status:  "error",
					Message: "Forbidden: insufficient permissions",
					Error:   "NO_PERMISSION",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthority(permission string, AuthorityPermissionCheck func(userId,restaurantId int64 ,targetRoleId int64, permission string) bool) func(http.Handler) http.Handler {
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

			if !AuthorityPermissionCheck(user.UserID,user.RestaurantId,roleID,permission) {
				utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
					Status:  "error",
					Message: "Forbidden: insufficient permissions",
					Error:   "NO_PERMISSION",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
