package middlewares

import (
	"net/http"

	middlewares "github.com/harshgupta9473/restaurantmanagement/helpers/middleware"
)

func IsVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims := r.Context().Value("user")
		if userClaims == nil {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		claims, ok := userClaims.(JWTClaims)
		if !ok {
			http.Error(w, "Invalid user claims format", http.StatusUnauthorized)
			return
		}
		if claims.Verified{
			next.ServeHTTP(w,r)
		}

		userID := claims.UserID

		err := middlewares.IsVerified(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
