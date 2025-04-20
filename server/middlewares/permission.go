package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/harshgupta9473/restaurantmanagement/utils"
)

// Just to check if  Permissions like "view_inventory etc and also some authority permission if they have"
func RequirePermission(permission string, checkPermissionFunc func(userID, restaurantID int64, permission string) (int64,error)) func(http.Handler) http.Handler {
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

			
			permissionId,err := checkPermissionFunc(user.UserID, user.RestaurantId, permission)
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

			ctx:=context.WithValue(r.Context(),"permission_id",permissionId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthority checks if the user has the required authority on a specific role.
func RequireAuthority(permission string, AuthorityPermissionCheck func(userId, restaurantId int64, targetRoleId int64, permission string) (int64,error)) func(http.Handler) http.Handler {
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

			
			permissionId, err := AuthorityPermissionCheck(user.UserID, user.RestaurantId, roleID, permission);
			if err != nil {
				
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

			ctx:=context.WithValue(r.Context(),"permission_id",permissionId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func GetPermissionIdFromContext(r *http.Request)(int64,error){
	id,ok:=r.Context().Value("permission_id").(int64)
	if !ok{
		return 0,fmt.Errorf("failed to access permission id from context")
	}
	return id,nil
}