package models

import "time"


type RestaurantFormReq struct {
	ID          int64  `json:"id"`
	OwnerID     int64  `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	StreetAddress string `json:"street_address"`
	Locality      string `json:"locality"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	FoodType  string  `json:"food_type"` // veg, non-veg, both

	ContactNumber string `json:"contact_number"`
	ContactEmail  string `json:"contact_email"`
	ImageURL      string `json:"image_url"`
	GSTNumber     string `json:"gst_number"`
	PANNumber     string `json:"pan_number"`

	Status string `json:"status"`

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Restaurant struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"owner_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StreetAddress string    `json:"street_address"`
	Locality      string    `json:"locality"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PostalCode    string    `json:"postal_code"`
	Country       string    `json:"country"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	FoodType      string    `json:"food_type"`
	ContactNumber string    `json:"contact_number"`
	ContactEmail  string    `json:"contact_email"`
	GSTNumber     string    `json:"gst_number"`
	PANNumber     string    `json:"pan_number"`
	IsActive      bool      `json:"is_active"`
	ImageURL      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RestaurantPublic struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"-"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StreetAddress string    `json:"street_address"`
	Locality      string    `json:"locality"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	PostalCode    string    `json:"postal_code"`
	Country       string    `json:"country"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	FoodType      string    `json:"food_type"`
	ContactNumber string    `json:"-"`
	ContactEmail  string    `json:"-"`
	GSTNumber     string    `json:"-"`
	PANNumber     string    `json:"-"`
	IsActive      bool      `json:"is_active"`
	ImageURL      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PublicRestaurantDetails struct {
	Restaurant RestaurantPublic  `json:"restaurant"`
	Details    RestaurantDetails `json:"details"`
	Review     int64             `json:"reviews"`
	Tags       []RestaurantTag   `json:"tags"`
}

type RestaurantDetails struct {
	ID            int64  `json:"id"`
	RestaurantID  int64  `json:"restaurant_id"`
	AvgCostForTwo int64  `json:"avg_cost_for_two"`
	OpeningTime   string `json:"opening_time"` // "15:00:00"
	ClosingTime   string `json:"closing_time"`
}

type RestaurantImage struct {
	ID           int64     `json:"id"`
	RestaurantID int64     `json:"restaurant_id"`
	ImageURL     string    `json:"image_url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type RestaurantReview struct {
	ID           int64     `json:"id"`
	RestaurantID int64     `json:"restaurant_id"`
	UserID       int64     `json:"user_id"`
	Rating       int64     `json:"rating"`
	Review       string    `json:"review"`
	CreatedAt    time.Time `json:"created_at"`
}

type RestaurantTag struct {
	ID           int64  `json:"id"`
	RestaurantID int64  `json:"restaurant_id"`
	Tag          string `json:"tag"`
}

type PrivateRestaurantDetails struct{
	Restaurant Restaurant
	Details RestaurantDetails
	Images []RestaurantImage
	Review int64
	Tags []RestaurantTag
}
