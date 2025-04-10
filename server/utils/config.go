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
