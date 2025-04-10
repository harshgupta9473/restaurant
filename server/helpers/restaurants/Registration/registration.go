package restaurants

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func CreateRestaurantRequest(req *models.RestaurantFormReq) (int64, error) {
	db := db.GetDB()
	query := `
		INSERT INTO restaurants (
			owner_id, title, street_address, locality, city, state,
			postal_code, country, latitude, longitude, food_type,
			contact_number, contact_email, image_url,
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14
		)
		RETURNING id
	`

	var restaurantID int64
	err := db.QueryRow(
		query,
		req.OwnerID,
		req.Title,
		req.StreetAddress,
		req.Locality,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.Latitude,
		req.Longitude,
		req.FoodType,
		req.ContactNumber,
		req.ContactEmail,
		req.ImageURL,
	).Scan(&restaurantID)

	return restaurantID, err
}

func BlockRestaurantRequest(restaurantID int64) (sql.Result,error) {
	db := db.GetDB()
	query := `
		UPDATE restaurantsreq
		SET status = 'block', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status != 'block'
	`
	res, err := db.Exec(query, restaurantID)
	return res, err
}

func DeletePendingRestaurantRequest(restaurantID int64) (sql.Result,error) {
	db := db.GetDB()
	query := `
		DELETE FROM restaurantsreq
		WHERE id = $1 and status='pending'
	`
	res, err := db.Exec(query, restaurantID)
	return res,err
}

func DeleteBlockedRestaurantRequest(restaurantID int64) (sql.Result,error) {
	db := db.GetDB()
	query := `
		DELETE FROM restaurantsreq
		WHERE id = $1 and status='block'
	`
	res, err := db.Exec(query, restaurantID)
	return res,err
}

func ApproveAndCreateRestaurant(requestID int64) (*models.Restaurant, error) {
	conn := db.GetDB()

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	var req models.RestaurantFormReq
	fetchQuery := `
		SELECT id, owner_id, title, street_address, locality, city, state,
		       postal_code, country, latitude, longitude, food_type,
		       contact_number, contact_email, image_url, gst_number, pan_number
		FROM restaurantsreq
		WHERE id = $1 AND status = 'pending'
	`
	err = tx.QueryRow(fetchQuery, requestID).Scan(
		&req.ID,
		&req.OwnerID,
		&req.Title,
		&req.StreetAddress,
		&req.Locality,
		&req.City,
		&req.State,
		&req.PostalCode,
		&req.Country,
		&req.Latitude,
		&req.Longitude,
		&req.FoodType,
		&req.ContactNumber,
		&req.ContactEmail,
		&req.ImageURL,
		&req.GSTNumber,
		&req.PANNumber,
	)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch request: %v", err)
	}


	approveQuery := `
		UPDATE restaurantsreq
		SET status = 'approve', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err = tx.Exec(approveQuery, requestID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	
	insertQuery := `
		INSERT INTO restaurants (
			owner_id, title, street_address, locality, city, state,
			postal_code, country, latitude, longitude, food_type,
			contact_number, contact_email, image_url, gst_number, pan_number
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16
		)
		RETURNING id
	`

	var restaurantID int64
	err = tx.QueryRow(
		insertQuery,
		req.OwnerID,
		req.Title,
		req.StreetAddress,
		req.Locality,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.Latitude,
		req.Longitude,
		req.FoodType,
		req.ContactNumber,
		req.ContactEmail,
		req.ImageURL,
		req.GSTNumber,
		req.PANNumber,
	).Scan(&restaurantID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	restaurant := &models.Restaurant{
		ID:            restaurantID,
		OwnerID:       req.OwnerID,
		Title:         req.Title,
		StreetAddress: req.StreetAddress,
		Locality:      req.Locality,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		FoodType:      req.FoodType,
		ContactNumber: req.ContactNumber,
		ContactEmail:  req.ContactEmail,
		ImageURL:      req.ImageURL,
		GSTNumber:     req.GSTNumber,
		PANNumber:     req.PANNumber,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return restaurant, nil
}
