package middlewares

import (
	"net/http"

	"github.com/harshgupta9473/restaurantmanagement/db"
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

		userID := claims.UserID

		db := db.GetDB()
		var verified bool
		err := db.QueryRow(`SELECT verified FROM users WHERE id=$1`, userID).Scan(&verified)
		if err != nil {
			http.Error(w, "Failed to fetch user verification status", http.StatusInternalServerError)
			return
		}

		if !verified {
			http.Error(w, "Please verify your account to proceed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
