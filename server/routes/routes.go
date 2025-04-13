package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/restaurantmanagement/handlers/userAuth"
	restaurants "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/Registration"
	 details "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/RestaurantDetails"
	admin "github.com/harshgupta9473/restaurantmanagement/handlers/superAdmin"
	"github.com/harshgupta9473/restaurantmanagement/middlewares"
)

func SetupRoutes() mux.Router {
	router := mux.NewRouter()

	//user auth
	userAuth := router.PathPrefix("/auth/user").Subrouter()
	userAuth.HandleFunc("/signup", auth.Signup).Methods(http.MethodPost)
	userAuth.HandleFunc("/login", auth.UserLogin).Methods(http.MethodPost)
	userAuth.HandleFunc("/verify", auth.Verify).Methods(http.MethodPost)
	userAuth.Handle("/verify", middlewares.AuthMiddleware(http.HandlerFunc((auth.SendVerificationLink)))).Methods(http.MethodGet)

	// resturant registration
	register := router.PathPrefix("/restaurant/register").Subrouter()
	register.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(restaurants.RestaurantAccountRequest))).Methods(http.MethodPost)

	//SuperAdmin

	// restaurant request
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Handle("/approve-restaurant/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.ApproveRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/block/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.BlockRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/request-again/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.ReRequestRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/delete/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.DeleteRestaurantRequest)))).Methods(http.MethodGet)

	
	   
	// restaurant details

	//public
	detailsRouter:=router.PathPrefix("/restaurant").Subrouter();
	detailsRouter.HandleFunc("/about/{id}",details.AboutRestaurant).Methods(http.MethodGet)

	// owner and admin
	detailsRouter.Handle("/about/{id}",middlewares.AuthMiddleware(middlewares.IsAllowedRolesForAboutRestaurantDetail(http.HandlerFunc(details.GetRestaurantsPrivateDetails))))


	//  /about/reviews/12?page=2&limit=5
	   detailsRouter.HandleFunc("/about/reviews/{id}",details.GetReview).Methods(http.MethodGet)



	   


	return *router
}
