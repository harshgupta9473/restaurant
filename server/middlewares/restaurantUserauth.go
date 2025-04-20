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

type JWTRestaurantClaims struct {
	UserID       int64 `json:"user_id"`
	RestaurantId int64 `json:"restaurant_id"`
	jwt.RegisteredClaims
}

func GenerateRestaurantAccessToken(userID int64, restaurantId int64) (string, error) {
	claims := JWTRestaurantClaims{
		UserID:       userID,
		RestaurantId: restaurantId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(utils.AccessRestaurantRoleJWTSecret)
}

func GenerateRestaurantRefreshToken(userID int64, restaurantId int64) (string, error) {
	claims := JWTRestaurantClaims{
		UserID:       userID,
		RestaurantId: restaurantId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(utils.RefreshRestaurantRoleJWTSecret)
}

func AuthRestaurantRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("restaurant_access_token")
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: no access token",
				Error:   "NO_ACCESS_TOKEN",
			})
			return
		}

		accessToken := cookie.Value
		claims, token, err := ValidateRestaurantRolesJWT(accessToken, utils.AccessRestaurantRoleJWTSecret)

		if err == nil && token.Valid {
			ctx := context.WithValue(r.Context(), "restaurant_user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		var verr *jwt.ValidationError
		if errors.As(err, &verr) && verr.Errors&jwt.ValidationErrorExpired != 0 {
			refreshCookie, err := r.Cookie("restaurant_refresh_token")
			if err != nil {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "No refresh token provided",
					Error:   "TOKEN_EXPIRED_NO_REFRESH",
				})
				return
			}

			refreshToken := refreshCookie.Value
			refreshClaims, rt, err := ValidateRestaurantRolesJWT(refreshToken, utils.RefreshRestaurantRoleJWTSecret)
			if err != nil || !rt.Valid {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "Invalid refresh token",
					Error:   "INVALID_REFRESH_TOKEN",
				})
				return
			}

			newAccess, err :=GenerateRestaurantAccessToken(refreshClaims.UserID, refreshClaims.RestaurantId)
			if err != nil {
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Failed to generate new access token",
					Error:   "ACCESS_TOKEN_GENERATION_FAILED",
				})
				return
			}
			newRefresh, err := GenerateRestaurantRefreshToken(refreshClaims.UserID, refreshClaims.RestaurantId)
			if err != nil {
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Failed to generate new refresh token",
					Error:   "REFRESH_TOKEN_GENERATION_FAILED",
				})
				return
			}

			utils.SetCookie(w, "restaurant_access_token", newAccess, 15*60)
			utils.SetCookie(w, "restaurant_refresh_token", newRefresh, 7*24*3600)

			ctx := context.WithValue(r.Context(), "restaurant_user", refreshClaims)
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

func ValidateRestaurantRolesJWT(tokenStr string, jwtSecret []byte) (*JWTRestaurantClaims, *jwt.Token, error) {
	claims := &JWTRestaurantClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, token, err
	}

	if !token.Valid {
		return nil, token, errors.New("invalid token")
	}

	return claims, token, nil
}

func GetRestaurantUserContext(r *http.Request) (*JWTRestaurantClaims, error) {
	user, ok := r.Context().Value("restaurant_user").(JWTRestaurantClaims)
	if !ok {
		return nil, fmt.Errorf("failed to access user claims from context")
	}
	return &user, nil
}
