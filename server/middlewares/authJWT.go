package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/harshgupta9473/restaurantmanagement/utils"
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Verified bool   `json:"verified"`
	jwt.RegisteredClaims
}



func GenerateTokenForRole(userID int64, verified bool) (string, string, error) {
	accessClaims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		Verified: verified,
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
		claims,token,err := ValidateJWT(tokenStr,string(utils.AccessJWTSecret))

		if err == nil && token.Valid {
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

			newClaims ,rt,err:= ValidateJWT(refreshToken,string(utils.RefreshJWTSecret))
			if err != nil || !rt.Valid {
				utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
					Status:  "error",
					Message: "Invalid refresh token",
					Error:   "INVALID_REFRESH_TOKEN",
				})
				return
			}

			// Generate new tokens
			newAccessToken, newRefreshToken, err := GenerateTokenForRole(newClaims.UserID, newClaims.Verified)
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

func ValidateJWT(tokenStr string,JWTSecret string)(*JWTClaims,*jwt.Token,error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return nil,nil, err
	}

	if !token.Valid {
		return nil, nil,errors.New("invalid token")
	}

	return claims,token, nil
}
