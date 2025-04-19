package roleModels

import "time"

type NewRoleCreationRequestOwner struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantId int64  `json:"restaurant_id"`
	Level        int64  `json:"level"`
}

type NewRoleCreationRequestPower struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Level        int64  `json:"level"`
}

type Roles struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"-"`
	Description  string    `json:"description"`
	RestaurantId int64     `json:"restaurant_id"`
	ManagerId    int64     `json:"manager_id"`
	Level        int64     `json:"level"`
	Is_Global    bool      `json:"is_global"`
	Is_Custom    bool      `json:"is_custom"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NewUserRequestForRole struct {
	RestaurantID int64  `json:"restaurant_id"`
	Code         string `json:"code"`
	RoleId       int64  `json:"role_id"`
}
