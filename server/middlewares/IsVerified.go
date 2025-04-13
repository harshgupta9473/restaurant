package middlewares

import (
	"net/http"

	middlewaresHelper "github.com/harshgupta9473/restaurantmanagement/helpers/middleware"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func IsVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userClaims := r.Context().Value("user")
		if userClaims == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "User not found in context",
				Error:   "Unauthorized",
			})
			return
		}

		claims, ok := userClaims.(JWTClaims)
		if !ok {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Invalid user claims format",
				Error:   "Unauthorized",
			})
			return
		}

		if claims.Verified {
			next.ServeHTTP(w, r)
			return
		}

		err := middlewaresHelper.IsVerified(claims.UserID)
		if err != nil {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "User is not verified",
				Error:   err.Error(),
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
