package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	restaurants "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/Registration"
	details "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/RestaurantDetails"
	roles "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/Roles"
	restaurantsAuth "github.com/harshgupta9473/restaurantmanagement/handlers/restaurants/restaurantAuth"
	admin "github.com/harshgupta9473/restaurantmanagement/handlers/superAdmin"
   auth	"github.com/harshgupta9473/restaurantmanagement/handlers/userAuth"
	middlewaresHelper "github.com/harshgupta9473/restaurantmanagement/helpers/middleware"
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

	   // super admin login
	  router.HandleFunc("/auth/admin/login",admin.SuperAdminLogin)

	   // restaurant owner login
       router.HandleFunc("/restaurant/owner/login",restaurantsAuth.RestaurantOwnerLogIN)




	//resturant registration Request by User
	register := router.PathPrefix("/restaurant/register").Subrouter()
	register.Handle("/", middlewares.AuthMiddleware(middlewares.IsVerified(http.HandlerFunc(restaurants.RestaurantAccountRequest)))).Methods(http.MethodPost)
	




	//SuperAdmin

	// LogIN

	router.HandleFunc("/admin/auth/login",admin.SuperAdminLogin)

	// restaurantApproval request
	adminRouter := router.PathPrefix("/admin/restaurants").Subrouter()
	adminRouter.Handle("/",middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.GetAllListOFRequest))))
	adminRouter.Handle("/approve-restaurant/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.ApproveRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/block/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.BlockRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/request-again/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.ReRequestRestaurantRequest)))).Methods(http.MethodGet)
	adminRouter.Handle("/delete/{id}", middlewares.AuthMiddleware(middlewares.IsSuperAdminMiddleware(http.HandlerFunc(admin.DeleteRestaurantRequest)))).Methods(http.MethodGet)

	

	
	// restaurant details

	    //public
	    detailsRouter:=router.PathPrefix("/restaurant").Subrouter();
	    detailsRouter.HandleFunc("/about/{id}",details.AboutRestaurant).Methods(http.MethodGet)
	    // /about/reviews/12?page=2&limit=5
	    detailsRouter.HandleFunc("/about/reviews/{id}",details.GetReview).Methods(http.MethodGet)

	    // owner and admin
	    detailsRouter.Handle("/about/{id}",middlewares.AuthMiddleware(middlewares.IsAllowedRolesForAboutRestaurantDetail(http.HandlerFunc(details.GetRestaurantsPrivateDetails))))


	// role

	// role login 
	 
	rolesRouter:=router.PathPrefix("/restaurant/roles").Subrouter();

	// role login 
	rolesRouter.HandleFunc("/auth/login",roles.RoleLogin)
	    
	// role request
	     // role requesting  by normal user
		 rolesRouter.Handle("/request/user",middlewares.AuthMiddleware(middlewares.IsVerified(http.HandlerFunc(roles.RequestRole)))).Methods("POST")

		 // roles under different authority
	     //manage
		 // all roles under authority
		 rolesRouter.Handle("/manage",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequirePermission("manage",middlewaresHelper.CheckPermissionFunc)(http.HandlerFunc(roles.GetAllRolesThatsUnderAuthority))))
		 // get all requests of role_id
		 rolesRouter.Handle("/manage/{roleID}",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequireAuthority("manage",middlewaresHelper.AuthorityPermissionCheck)(http.HandlerFunc(roles.GetAllRequestInRoleThatsUnderAuthority))))

		 rolesRouter.Handle("/manage/approve/{roleID}/{staffID}",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequireAuthority("manage",middlewaresHelper.AuthorityPermissionCheck)(http.HandlerFunc(roles.ApproveRoleRequest))))
		 rolesRouter.Handle("/manage/reject/{roleID}/{staffID}",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequireAuthority("manage",middlewaresHelper.AuthorityPermissionCheck)(http.HandlerFunc(roles.ApproveRoleRequest))))
		 rolesRouter.Handle("/manage/block/{roleID}/{staffID}",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequireAuthority("manage",middlewaresHelper.AuthorityPermissionCheck)(http.HandlerFunc(roles.ApproveRoleRequest))))

		 // create roles by restaurant owner

		 router.Handle("/restaurant/owner/roles/create",middlewares.AuthMiddleware(middlewares.IsVerified(http.HandlerFunc(roles.CreateRolesByOwner))))

		 // create roles by those who have power to create roles
		 rolesRouter.Handle("/create/new",middlewares.AuthRestaurantRoleMiddleware(middlewares.RequirePermission("create_roles",middlewaresHelper.CheckPermissionFunc)(http.HandlerFunc(roles.CreateRolesWhoHavePowerToCreate))))




		 //menu


	   


	return *router
}
