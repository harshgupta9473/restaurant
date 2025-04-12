package restaurants

import (
	"database/sql"

	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func GetCompletePublicRestaurantDetails(restaurantID int64) (*models.PublicRestaurantDetails, error) {
	var result models.PublicRestaurantDetails

	db := db.GetDB()

	queryRestaurant := `SELECT id, owner_id, title, description, street_address, locality, city, state, postal_code, country, latitude, longitude, food_type, is_active, image_url, created_at, updated_at 
	                    FROM restaurants WHERE id = $1`
	err := db.QueryRow(queryRestaurant, restaurantID).Scan(
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
		&result.Restaurant.IsActive,
		&result.Restaurant.ImageURL,
		&result.Restaurant.CreatedAt,
		&result.Restaurant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryDetails := `SELECT id, restaurant_id, avg_cost_for_two, opening_time, closing_time FROM restaurant_details WHERE restaurant_id = $1`
	err = db.QueryRow(queryDetails, restaurantID).Scan(
		&result.Details.ID,
		&result.Details.RestaurantID,
		&result.Details.AvgCostForTwo,
		&result.Details.OpeningTime,
		&result.Details.ClosingTime,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	queryReviewsNo := `SELECT count(user_id) FROM restaurant_reviews WHERE restaurant_id = $1`
	err = db.QueryRow(queryReviewsNo, restaurantID).Scan(&result.Review)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	queryTags := `SELECT id, restaurant_id, tag FROM restaurant_tags WHERE restaurant_id = $1`
	tagRows, err := db.Query(queryTags, restaurantID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tag models.RestaurantTag
		err := tagRows.Scan(&tag.ID, &tag.RestaurantID, &tag.Tag)
		if err != nil {
			return nil, err
		}
		result.Tags = append(result.Tags, tag)
	}

	return &result, nil
}

