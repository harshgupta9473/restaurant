package utils

import (
	"os"
)

var (
	AccessJWTSecret  []byte
	RefreshJWTSecret []byte
	AccessRestaurantRoleJWTSecret []byte
	RefreshRestaurantRoleJWTSecret []byte
)

func LoadSecrets() {
	AccessJWTSecret = []byte(os.Getenv("accessJWTSecretKey"))
	RefreshJWTSecret = []byte(os.Getenv("refreshJWTSecretKey"))
	AccessRestaurantRoleJWTSecret=[]byte(os.Getenv("accessRestaurantRoleJWTSecretKey"))
	RefreshRestaurantRoleJWTSecret=[]byte(os.Getenv("refreshRestaurantRoleJWTSecretKey"))
}

var (
	AllowedRolesForAll=[]string{"user","admin","owner","waiter","inventory","staff"}
	AllowedRolesForAboutRestaurantDetails=[]string{"admin","owner"}
)

