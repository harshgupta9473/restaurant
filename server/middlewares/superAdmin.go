package middlewares

import (
	"net/http"

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

		if user.Role!="admin" {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Access denied: Not a super admin",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
