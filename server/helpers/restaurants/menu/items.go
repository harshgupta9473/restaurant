package menu

import (
	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func AddMenuItem(item models.MenuItem) error {
	conn := db.GetDB()

	query := `
		INSERT INTO menu_items (restaurant_id, category_id, name, description, price, cuisine, image_url, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := conn.Exec(query,
		item.RestaurantId,
		item.CategoryId,
		item.Name,
		item.Description,
		item.Price,
		item.Cuisine,
		item.ImageURL,
		item.IsAvailable,
	)
	return err
}

func UpdateMenuItem(item models.MenuItem) error {
	conn := db.GetDB()

	query := `
		UPDATE menu_items
		SET category_id = $1,
		    name = $2,
		    description = $3,
		    price = $4,
		    cuisine = $5,
		    image_url = $6,
		    is_available = $7,
		    updated_at = NOW()
		WHERE id = $8 AND restaurant_id = $9
	`
	_, err := conn.Exec(query,
		item.CategoryId,
		item.Name,
		item.Description,
		item.Price,
		item.Cuisine,
		item.ImageURL,
		item.IsAvailable,
		item.ID,
		item.RestaurantId,
	)

	return err
}

func DeleteMenuItem(itemID int64, restaurantID int64) error {
	conn := db.GetDB()

	query := `
		DELETE FROM menu_items
		WHERE id = $1 AND restaurant_id = $2
	`
	_, err := conn.Exec(query, itemID, restaurantID)
	return err
}

func MakeMenuItemUnavailable(itemID int64, restaurantID int64) error {
	conn := db.GetDB()

	query := `
		UPDATE menu_items
		SET is_available = FALSE, updated_at = NOW()
		WHERE id = $1 AND restaurant_id = $2
	`
	_, err := conn.Exec(query, itemID, restaurantID)
	return err
}


func GetAllMenuItemsForRestaurant(restaurantId int64) ([]models.MenuItem, error) {
	conn := db.GetDB()

	query := `
		SELECT id, restaurant_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		FROM menu_items
		WHERE restaurant_id = $1
	`
	rows, err := conn.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem

	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(
			&item.ID,
			&item.RestaurantId,
			&item.CategoryId,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.ImageURL,
			&item.IsAvailable,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
