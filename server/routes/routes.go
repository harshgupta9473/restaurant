package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/restaurantmanagement/handlers/auth"
	restaurants "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/Registration"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
)

func SetupRoutes() mux.Router {
	router:=mux.NewRouter();

	//user auth
	userAuth:=router.PathPrefix("/auth/user").Subrouter()
	userAuth.HandleFunc("/signup",auth.Signup).Methods(http.MethodPost)
	userAuth.HandleFunc("/login",auth.Login).Methods(http.MethodPost)
	userAuth.HandleFunc("/verify",auth.Verify).Methods(http.MethodPost)
	userAuth.Handle("/verify",middlewares.AuthMiddleware(http.HandlerFunc((auth.SendVerificationLink)))).Methods(http.MethodGet)


	// resturant registration
	register:=router.PathPrefix("restaurant/register").Subrouter()
	register.Handle("/",middlewares.AuthMiddleware(http.HandlerFunc(restaurants.RestaurantAccountRequest))).Methods(http.MethodPost)
	

	



	//SuperAdmin 

	    // restaurant request
	    admin:=router.PathPrefix("/admin").Subrouter()
	    admin.Handle("/approve-restaurant/{id}",middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(restaurants.ApproveRestaurantRequest)))).Methods(http.MethodGet)
	    admin.Handle("/block/{id}",middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(restaurants.BlockRestaurantRequest)))).Methods(http.MethodGet)
	    admin.Handle("/request-again/{id}",middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(restaurants.ReRequestRestaurantRequest)))).Methods(http.MethodGet)
	    admin.Handle("/request-again/{id}",middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(restaurants.DeleteRestaurantRequest)))).Methods(http.MethodGet)

		//

	return *router
}