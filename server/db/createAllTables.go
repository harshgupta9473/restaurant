package db

import (
	"fmt"
	"log"
)

func CreateAllTable() error {
	err := CreateUserTable()
	if err != nil {
		return fmt.Errorf("error creating user table: %w", err)
	}
	log.Println("user table created")

	err = CreateOTPTable()
	if err != nil {
		return fmt.Errorf("error creating OTP table: %w", err)
	}
	log.Println("OTP table created")

	err = CreateRestaurantReqTable()
	if err != nil {
		return fmt.Errorf("error creating restaurant request table: %w", err)
	}
	log.Println("restaurant request table created")

	err = CreateSuperAdmin()
	if err != nil {
		return fmt.Errorf("error creating super admin table: %w", err)
	}
	log.Println("super admin table created")

	err = CreateRestaurantsTable()
	if err != nil {
		return fmt.Errorf("error creating restaurants table: %w", err)
	}
	log.Println("restaurants table created")

	err = CreateRestaurantDetailsTable()
	if err != nil {
		return fmt.Errorf("error creating restaurant details table: %w", err)
	}
	log.Println("restaurant details table created")

	err = CreateRestaurantImagesTable()
	if err != nil {
		return fmt.Errorf("error creating restaurant images table: %w", err)
	}
	log.Println("restaurant images table created")

	err = CreateRestaurantReviewsTable()
	if err != nil {
		return fmt.Errorf("error creating restaurant reviews table: %w", err)
	}
	log.Println("restaurant reviews table created")

	err = CreateRestaurantTagsTable()
	if err != nil {
		return fmt.Errorf("error creating restaurant tags table: %w", err)
	}
	log.Println("restaurant tags table created")

	err = CreateRolesTable()
	if err != nil {
		return fmt.Errorf("error creating roles table: %w", err)
	}
	log.Println("roles table created")

	err = CreatePermissionsTable()
	if err != nil {
		return fmt.Errorf("error creating permissions table: %w", err)
	}
	log.Println("permissions table created")

	err = CreateRolePermissionsTable()
	if err != nil {
		return fmt.Errorf("error creating role_permissions table: %w", err)
	}
	log.Println("role_permissions table created")

	err = CreateRolesTargetPermission()
	if err != nil {
		return fmt.Errorf("error creating roles_target_permissions table: %w", err)
	}
	log.Println("roles_target_permissions table created")

	err = CreateStaffRoles()
	if err != nil {
		return fmt.Errorf("error creating staff_roles table: %w", err)
	}
	log.Println("staff_roles table created")

	return nil
}
