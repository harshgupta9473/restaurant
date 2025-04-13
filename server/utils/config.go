package utils

import (
	"os"
)

var (
	AccessJWTSecret  []byte
	RefreshJWTSecret []byte
)

func LoadSecrets() {
	AccessJWTSecret = []byte(os.Getenv("accessJWTSecretKey"))
	RefreshJWTSecret = []byte(os.Getenv("refreshJWTSecretKey"))
}

var (
	AllowedRolesForAll=[]string{"user","admin","owner","waiter","inventory","staff"}
	AllowedRolesForAboutRestaurantDetails=[]string{"admin","owner"}
)

