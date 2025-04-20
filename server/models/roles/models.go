package roleModels

import "time"

type RestaurantRoleLogin struct {
	Email        string `json:"email"`
	RestaurantId int64  `json:"restaurant_id"`
	Password     string `json:"password"`
}

type NewRoleCreationRequestOwner struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantId int64  `json:"restaurant_id"`
	Level        int64  `json:"level"`
}

type NewRoleCreationRequestPower struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
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

type StaffMember struct {
	ID int64 `json:"id"`
}
type PersonWithRoles struct {
	StaffId      int64     `json:"staff_id"`
	UserID       int64     `json:"user_id"`
	RestaurantId int64     `json:"restaurant_id"`
	RoleId       int64     `json:"role_id"`
	RoleName     string    `json:"role_name"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	MiddleName   string    `json:"middle_name"`
	ManagerId    int64     `json:"manager_id"`
	Blocked      bool      `json:"blocked"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Is_Approved  int       `json:"is_approved"` // 0 initital 1 approved -1 rejected -2 blocked
}
