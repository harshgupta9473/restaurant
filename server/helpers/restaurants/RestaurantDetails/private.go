package restaurants

import (
	"database/sql"

	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func GetDetailsOfRestaurant(restaurantId int64) (*models.PrivateRestaurantDetails, error) {
	conn := db.GetDB()
	var result models.PrivateRestaurantDetails

	
	queryRestaurant := `SELECT * FROM restaurants WHERE id=$1`
	err := conn.QueryRow(queryRestaurant, restaurantId).Scan(
		&result.Restaurant.ID,
		&result.Restaurant.OwnerID,
		&result.Restaurant.Title,
		&result.Restaurant.Description,
		&result.Restaurant.StreetAddress,
		&result.Restaurant.Locality,
		&result.Restaurant.City,
		&result.Restaurant.State,
		&result.Restaurant.PostalCode,
		&result.Restaurant.Country,
		&result.Restaurant.Latitude,
		&result.Restaurant.Longitude,
		&result.Restaurant.FoodType,
		&result.Restaurant.ContactNumber,
		&result.Restaurant.ContactEmail,
		&result.Restaurant.GSTNumber,
		&result.Restaurant.PANNumber,
		&result.Restaurant.IsActive,
		&result.Restaurant.ImageURL,
		&result.Restaurant.CreatedAt,
		&result.Restaurant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryDetails := `SELECT * FROM restaurant_details WHERE restaurant_id=$1`
	err = conn.QueryRow(queryDetails, restaurantId).Scan(
		&result.Details.ID,
		&result.Details.RestaurantID,
		&result.Details.AvgCostForTwo,
		&result.Details.OpeningTime,
		&result.Details.ClosingTime,
	)
	if err != nil && err != sql.ErrNoRows {
		return &result, err
	}

	queryImages := `SELECT * FROM restaurant_images WHERE restaurant_id=$1`
	rows, err := conn.Query(queryImages, restaurantId)
	if err != nil && err != sql.ErrNoRows {
		return &result, err
	}
	defer rows.Close()
	for rows.Next() {
		var img models.RestaurantImage
		if err := rows.Scan(
			&img.ID,
			&img.RestaurantID,
			&img.ImageURL,
			&img.UploadedAt,
		); err != nil {
			return &result, err
		}
		result.Images = append(result.Images, img)
	}

	
	queryReview := `SELECT COUNT(id) FROM restaurant_reviews WHERE restaurant_id=$1`
	err = conn.QueryRow(queryReview, restaurantId).Scan(&result.Review)
	if err != nil && err != sql.ErrNoRows {
		return &result, err
	}

	
	queryTags := `SELECT * FROM restaurant_tags WHERE restaurant_id=$1`
	tagRows, err := conn.Query(queryTags, restaurantId)
	if err != nil && err != sql.ErrNoRows {
		return &result, err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var tag models.RestaurantTag
		if err := tagRows.Scan(
			&tag.ID,
			&tag.RestaurantID,
			&tag.Tag,
		); err != nil {
			return &result, err
		}
		result.Tags = append(result.Tags, tag)
	}

	return &result, nil
}


