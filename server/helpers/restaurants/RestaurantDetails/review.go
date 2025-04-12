package restaurants

import (
	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func GetReview(restaurantId int64, page, limit int) ([]models.RestaurantReview, error) {
	db := db.GetDB()
	var reviews []models.RestaurantReview

	offset := (page - 1) * limit

	query := `SELECT id, restaurant_id, user_id, rating, review, created_at 
	          FROM restaurant_reviews 
	          WHERE restaurant_id = $1 
	          ORDER BY created_at DESC 
	          LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, restaurantId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var review models.RestaurantReview
		err := rows.Scan(&review.ID, &review.RestaurantID, &review.UserID, &review.Rating, &review.Review, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func AddReview(userId int64, restaurantId int64, text string, rating int) (*models.RestaurantReview, error) {
	conn := db.GetDB()
	query := `Insert into restaurant_reviews
	(restaurant_id,user_id,rating,review)
	values($1,$2,$3,$4)
	returning *
	`
	var res models.RestaurantReview
	err := conn.QueryRow(query, restaurantId, userId, rating, text).Scan(
		&res.ID,
		&res.RestaurantID,
		&res.UserID,
		&res.Rating,
		&res.Review,
		&res.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
