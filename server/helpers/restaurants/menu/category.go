package menu

import (
	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/restaurants"
)

func CreateMenuCategory(restaurantID *int64, name string, isCustom bool)error {
	conn := db.GetDB()
	query:=`INSERT INTO menu_categories(restaurant_id,name,is_custom)
	VALUES($1,$2,$3)`
	_,err:=conn.Exec(query,restaurantID,name,isCustom)
	return err
}

// func DeleteCategoryByAdmin(category_id int64){

// }

func DeleteCategoryByRestaurant(restaurant_id *int64,category_id int64)error{
	conn:=db.GetDB()
	query:=`DELETE FROM menu_categories where id=$1 and restaurant_id=$2`
	_,err:=conn.Exec(query,category_id,restaurant_id)
	return err
}


// func UpdateCategory(restaurant_id *int64, category_id int64, newName string, newIsCustom bool) (*models.MenuCategory, error) {
// 	conn := db.GetDB()


// 	query := `
// 		UPDATE menu_categories
// 		SET name = $1, is_custom = $2, updated_at = NOW()
// 		WHERE id = $3 AND restaurant_id = $4
// 		RETURNING id, name, restaurant_id, is_custom, created_at, updated_at
// 	`

	
// 	var updatedCategory models.MenuCategory

// 	err := conn.QueryRow(query, newName, newIsCustom, category_id, restaurant_id).Scan(
// 		&updatedCategory.ID,
// 		&updatedCategory.Name,
// 		&updatedCategory.RestaurantId,
// 		&updatedCategory.IsCustom,
// 		&updatedCategory.CreatedAt,
// 		&updatedCategory.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &updatedCategory, nil
// }



func GetAllCategoriesForRestaurant(restaurantId int64)([]models.MenuCategory,error){
	conn:=db.GetDB()
	query:=`SELECT id,name,restaurant_id,is_custom,created_at,updated_at FROM menu_categories WHERE restaurant_id is NULL OR restaurant_id=$1`
	rows,err:=conn.Query(query,restaurantId)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	var categories []models.MenuCategory

	for rows.Next(){
		var cat models.MenuCategory
		err:=rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.RestaurantId,
			&cat.IsCustom,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err!=nil{
			return nil,err
		}
		categories=append(categories, cat)
	}
	return categories,nil
}
