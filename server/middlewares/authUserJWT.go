package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Verified bool   `json:"verified"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokenForRole(userID int64, verified bool, role string) (string, string, error) {
	accessClaims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		Verified: verified,
		Role:     role,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString(utils.AccessJWTSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
		Verified: verified,
		Role:     role,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString(utils.RefreshJWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: no access token",
				Error:   "NO_ACCESS_TOKEN",
			})
			return
		}

		tokenStr := cookie.Value
		claims, token, err := ValidateJWT(tokenStr, utils.AccessJWTSecret)

		if err == nil && token.Valid && CheckAllowedRoles(claims.Role, utils.AllowedRolesForAll) {
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		//If the error returned by jwt.ParseWithClaims is of type *jwt.ValidationError and
		// the cause of the error was token expiration, then handle it accordingly.
		var verr *jwt.ValidationError
		if errors.As(err, &verr) && verr.Errors&jwt.ValidationErrorExpired != 0 {
			// Try refresh
			cookie, err := r.Cookie("refresh_token")
			if err != nil {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "No refresh token provided",
					Error:   "TOKEN_EXPIRED_NO_REFRESH",
				})
				return
			}
			refreshToken := cookie.Value
			if refreshToken == "" {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "No refresh token provided",
					Error:   "TOKEN_EXPIRED_NO_REFRESH",
				})
				return
			}

			newClaims, rt, err := ValidateJWT(refreshToken, utils.RefreshJWTSecret)
			if err != nil || !rt.Valid || !CheckAllowedRoles(newClaims.Role, utils.AllowedRolesForAll) {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "Invalid refresh token",
					Error:   "INVALID_REFRESH_TOKEN",
				})
				return
			}
			// Generate new tokens
			newAccessToken, newRefreshToken, err := GenerateTokenForRole(newClaims.UserID, newClaims.Verified, newClaims.Role)
			if err != nil {
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Failed to refresh tokens",
					Error:   "TOKEN_REFRESH_ERROR",
				})
				return
			}

			// set access token in cookie
			utils.SetCookie(w, "access_token", newAccessToken, 24*3600)

			// set Refresh token in cookie
			utils.SetCookie(w, "refresh_token", newRefreshToken, 7*24*3600)

			ctx := context.WithValue(r.Context(), "user", newClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid token",
			Error:   "INVALID_TOKEN",
		})
	})
}

func ValidateJWT(tokenStr string, JWTSecret []byte) (JWTClaims, *jwt.Token, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return *claims, nil, err
	}

	if !token.Valid {
		return *claims, nil, errors.New("invalid token")
	}

	return *claims, token, nil
}

func GetUserContext(r *http.Request) (*JWTClaims, error) {
	user, ok := r.Context().Value("user").(JWTClaims)
	if !ok {
		return nil, fmt.Errorf("failed to access user claims from context")
	}
	return &user, nil
}

func CheckAllowedRoles(userRole string, allowedroles []string) bool {
	for _, role := range allowedroles {
		if role == userRole {
			return true
		}
	}
	return false
}
