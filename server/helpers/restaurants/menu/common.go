package menu

import (
	"database/sql"

	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func GetAllMenuItemsWithCategory(restaurantId int64) ([]models.MenuItemsWithCategory, error) {
	conn := db.GetDB()
	query := `
		SELECT id, name, restaurant_id, is_custom, created_at, updated_at
		FROM menu_categories
		WHERE restaurant_id IS NULL OR restaurant_id = $1
	`

	rows, err := conn.Query(query, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MenuItemsWithCategory

	for rows.Next() {
		var cat models.MenuItemsWithCategory

		err := rows.Scan(
			&cat.MenuCategory.ID,
			&cat.MenuCategory.Name,
			&cat.MenuCategory.RestaurantId,
			&cat.MenuCategory.IsCustom,
			&cat.MenuCategory.CreatedAt,
			&cat.MenuCategory.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		menuItemQuery := `
			SELECT id, restaurant_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
			FROM menu_items
			WHERE category_id = $1 AND restaurant_id = $2
		`

		menuRows, err := conn.Query(menuItemQuery, cat.MenuCategory.ID, restaurantId)
		if err != nil {
			return nil, err
		}

		var items []models.MenuItem
		for menuRows.Next() {
			var item models.MenuItem
			err := menuRows.Scan(
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
				menuRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		menuRows.Close()

		cat.MenuItems = items
		result = append(result, cat)
	}

	return result, nil
}

func GetCategoryWithMenuItems(restaurantId int64, categoryId int64) (*models.MenuItemsWithCategory, error) {
	conn := db.GetDB()

	
	queryCategory := `
		SELECT id, name, restaurant_id, is_custom, created_at, updated_at
		FROM menu_categories
		WHERE id = $1 AND (restaurant_id IS NULL OR restaurant_id = $2)
	`

	var cat models.MenuItemsWithCategory

	err := conn.QueryRow(queryCategory, categoryId, restaurantId).Scan(
		&cat.MenuCategory.ID,
		&cat.MenuCategory.Name,
		&cat.MenuCategory.RestaurantId,
		&cat.MenuCategory.IsCustom,
		&cat.MenuCategory.CreatedAt,
		&cat.MenuCategory.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	
	queryItems := `
		SELECT id, restaurant_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		FROM menu_items
		WHERE category_id = $1 AND restaurant_id = $2
	`
	rows, err := conn.Query(queryItems, categoryId, restaurantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		cat.MenuItems = append(cat.MenuItems, item)
	}

	return &cat, nil
}
