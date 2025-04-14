package middlewares

import (
	"net/http"

	"github.com/harshgupta9473/restaurantmanagement/utils"
)

func IsSuperAdminMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		claim,err:=GetUserContext(r)
		if err!=nil{
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: User context not found",
				Error:   "NO_USER_CONTEXT",
			})
			return
		}
		if claim.Role!="admin"{
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Forbidden: Role not permitted for this route",
				Error:   "ROLE_NOT_ALLOWED",
			})
			return
		}
		next.ServeHTTP(w,r)
	})
}

