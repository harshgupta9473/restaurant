package middlewares

import (
	"net/http"

	middlewares "github.com/harshgupta9473/restaurantmanagement/helpers/middleware"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func IsSuperAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(JWTClaims)
		if !ok {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Invalid user context",
			})
			return
		}

		if user.UserID == 0 {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: Missing user ID",
			})
			return
		}

		isAdmin, err := middlewares.IsSuperAdmin(user.UserID)
		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Failed to verify super admin status",
				Error:   err.Error(),
			})
			return
		}

		if !isAdmin {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Access denied: Not a super admin",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
