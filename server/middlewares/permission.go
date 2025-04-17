package middlewares

import (
	"net/http"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)



func PermissionMiddleware(requiredPermission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims, err :=  GetRestaurantUserContext(r)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized",
				Error:   "NO_USER_CONTEXT",
			})
			return
		}

		if !HasPermission(userClaims.Permissions, requiredPermission) {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Forbidden: missing permission",
				Error:   "PERMISSION_DENIED",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HasPermission(userPermissions []string, required string) bool {
	for _, p := range userPermissions {
		if p == required {
			return true
		}
	}
	return false
}
