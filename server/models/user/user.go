package models

import (
	"time"
)

type User struct {
	Id         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Blocked    bool      `json:"blocked"`
}

type UserSignup struct {
	Id         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OTP struct {
	Id       int64     `json:"id"`
	UserId   int64     `json:"user_id"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
	Verified bool      `json:"verified"`
}

type Verify struct {
	UserId string `json:"user_id"` //user_id from user table
	Token  string `json:"token"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
