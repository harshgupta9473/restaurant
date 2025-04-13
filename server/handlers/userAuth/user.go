package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"
	"time"

	helpers "github.com/harshgupta9473/restaurantmanagement/helpers/auth"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
	models "github.com/harshgupta9473/restaurantmanagement/models/user"
	"github.com/harshgupta9473/restaurantmanagement/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.UserSignup
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid format",
			Error:   err.Error(),
		})
		return
	}

	_, err = helpers.GetUserByEmail(req.Email)
	if err == nil {
		utils.WriteJson(w, http.StatusConflict, utils.APIResponse{
			Status:  "error",
			Message: "User already exists",
			Error:   "USER_EXISTS",
		})
		return
	}

	if err != sql.ErrNoRows {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Internal Server Error : %w", err),
			Error:   "DB_ERROR",
		})
		return
	}

	token, err := utils.GenerateOTPToken()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Internal Server Error : %w", err),
			Error:   "ERR_OTP_GENERATE",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Internal Server Error : %w", err),
			Error:   "PASSWORD_HASH_ERROR",
		})
		return
	}

	req.Password = string(hashedPassword)
	userId, err := helpers.CreateNewUser(req)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Error creating user",
			Error:   err.Error(),
		})
		return
	}
	err = helpers.InsertTokenIntoOTPTable(token, userId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Error generating and inserting token in OTP table",
			Error:   err.Error(),
		})
		return
	}
	err = utils.SendVerificationEmail(req.Email, token, userId)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Error sending mail",
			Error:   err.Error(), // or use a custom tag like "MAIL_SEND_ERROR"
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Account created. Please verify your email to continue.",
	})
}

func Verify(w http.ResponseWriter, r *http.Request) {
	var req models.Verify
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Request",
			Error:   "Not the Correct Format",
		})
		return
	}

	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Internal Server Error",
			Data:    nil,
			Error:   "Error concerting the userid into int64",
		})
		return
	}

	details, err := helpers.FetchOTPDetailsByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: "Invalid Link",
				Error:   "user id is wrong",
			})
		} else {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   "DB Error",
			})
		}
		return
	}

	if details.Verified {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "User is already verified",
		})
		return
	}

	if details.Token != req.Token {
		utils.WriteJson(w, http.StatusBadGateway, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Request",
			Error:   "Wrong Token",
		})
		return
	}

	if details.ExpireAt.Before(time.Now()) {
		utils.WriteJson(w, http.StatusGatewayTimeout, utils.APIResponse{
			Status:  "error",
			Message: "TimedOut Request Verification Again",
			Error:   "Link Expired",
		})
		return
	}

	err = helpers.MarkUserVerified(userId)
	if err != nil {
		utils.WriteJson(w, http.StatusGatewayTimeout, utils.APIResponse{
			Status:  "error",
			Message: "DB ERROR",
			Error:   err.Error(),
		})
		return
	}

	//if user is already logged in?
	cookie, err := r.Cookie("access_token")
	if err == nil {
		claims, _, err := middlewares.ValidateJWT(cookie.Value, utils.AccessJWTSecret)
		if err == nil && claims.UserID == userId {
			newAccessToken, newRefreshToken, err := middlewares.GenerateTokenForRole(claims.UserID, true, claims.Role)
			if err == nil {
				// set access token in cookie
				utils.SetCookie(w, "access_token", newAccessToken, 24*3600)
				// set Refresh token in cookie
				utils.SetCookie(w, "refresh_token", newRefreshToken, 7*24*3600)
			}
		}
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "User is verified",
	})
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Request",
			Error:   "Not the Correct Format",
		})
		return
	}

	user,yes:=LoginHelper(w,r,req.Email,req.Password,"user")
	if !yes{
		return
	}


	accessToken, refreshToken, err := middlewares.GenerateTokenForRole(user.Id, user.Verified, "user")
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not generate token",
			Error:   err.Error(),
		})
		return
	}

	// set access token in cookie
	utils.SetCookie(w, "access_token", accessToken, 24*3600)

	// set Refresh token in cookie
	utils.SetCookie(w, "refresh_token", refreshToken, 7*24*3600)

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Login successful",
	})

}

func LoginHelper(w http.ResponseWriter, r *http.Request, email, password, role string) (*models.User,bool) {
	user, err := helpers.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: "Invalid Request",
				Error:   "Invalid Email or Password",
			})
		} else {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   "DB_ERROR",
			})
		}
		return nil, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Email Or Password",
			Error:   "Invalid_EMAIL_PASSWORD",
		})
		return nil, false
	}
	return user,true	
}

func SendVerificationLink(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("user").(middlewares.JWTClaims)
	if !ok {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized",
			Error:   "Invalid token claims",
		})
		return
	}

	if claims.Verified {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Already Verified",
			Error:   "User already verified",
		})
		return
	}

	userID := claims.UserID

	user, err := helpers.GetUserByUserId(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "User fetch failed",
			Error:   err.Error(),
		})
		return
	}

	token, err := utils.GenerateOTPToken()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Token generation failed",
			Error:   err.Error(),
		})
		return
	}

	err = helpers.InsertTokenIntoOTPTable(token, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to save OTP",
			Error:   err.Error(),
		})
		return
	}

	err = utils.SendVerificationEmail(user.Email, token, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Email Sending Failed",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Verification email sent",
	})
}
