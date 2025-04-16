package seed

import (
	"fmt"
	"log"
	"os"

	"github.com/harshgupta9473/restaurantmanagement/db"
)

func Seed() error {
	if os.Getenv("permission_seed") == "true" {
		err := SeedPermissions()
		if err != nil {
			return fmt.Errorf("error seeding permissions: %w", err)
		}
		log.Println("Permissions seeded successfully.")
	}

	if os.Getenv("role_seed") == "true" {
		err := SeedSystemRoles()
		if err != nil {
			return fmt.Errorf("error seeding roles: %w", err)
		}
		log.Println("Roles seeded successfully.")
	}

	if os.Getenv("role_permission_seed") == "true" {
		err := SeedRolePermission()
		if err != nil {
			return fmt.Errorf("error seeding role-permissions: %w", err)
		}
		log.Println("Role-permissions seeded successfully.")
	}

	return nil
}



func SeedPermissions() error {
	conn := db.GetDB()
	permissions := []struct {
		Name        string
		Description string
	}{
		// General permissions
		{"view_orders", "Can view all restaurant orders"},
		{"edit_menu", "Can add or update menu items"},
		{"view_reports", "Can view sales and performance reports"},
		{"manage_inventory", "Can update stock and inventory"},
		{"respond_reviews", "Can respond to customer reviews"},
		{"handle_reservations", "Can handle table bookings and reservations"},

		// Authority permissions
		{"assign_role", "Can assign roles to staff"},
		{"approve_staff", "Can approve new staff for the restaurant"},
		{"remove_staff", "Can remove staff from restaurant"},
		{"create_custom_role", "Can create new roles"},
		{"delete_role", "Can delete existing roles"},
		{"update_role_permissions", "Can update permissions for roles"},
	}

	query := `
		INSERT INTO permissions (name, description)
		VALUES ($1, $2)
		ON CONFLICT (name) DO NOTHING;
	`

	for _, p := range permissions {
		_, err := conn.Exec(query, p.Name, p.Description)
		if err != nil {
			return fmt.Errorf("failed to insert permission %s: %w", p.Name, err)
		}
	}

	return nil
}


func SeedSystemRoles() error {
	conn := db.GetDB()

	roles := []struct {
		Name        string
		Description string
		Level       int
	}{
		{"Manager", "Manages the overall operations of the restaurant", 3},
		{"Chef", "Prepares meals and manages kitchen", 2},
		{"Waiter", "Serves customers and handles orders", 1},
		{"Cleaner", "Responsible for cleanliness and hygiene", 0},
	}

	query := `
	INSERT INTO roles (name, description, level, is_global, is_custom)
	VALUES ($1, $2, $3, TRUE, FALSE)
	ON CONFLICT (name, restaurant_id) DO NOTHING;
	`

	for _, r := range roles {
		_, err := conn.Exec(query, r.Name, r.Description, r.Level)
		if err != nil {
			return fmt.Errorf("failed to insert system role %s: %w", r.Name, err)
		}
	}

	log.Println("System roles seeded successfully.")
	return nil
}


func SeedRolePermission()error{
	conn:=db.GetDB()
	rolePermission:=map[string][]string{
		"Manager": {
			"view_orders",
			"edit_menu",
			"view_reports",
			"manage_inventory",
			"respond_reviews",
			"handle_reservations",
			"assign_role",
			"approve_staff",
			"remove_staff",
			"create_custom_role",
			"delete_role",
			"update_role_permissions",
		},
		"Chef": {
			"edit_menu",
			"manage_inventory",
		},
		"Waiter": {
			"view_orders",
			"handle_reservations",
		},
	}
	for roleName,permissionNames:=range rolePermission{
		var roleId int64
		err:=conn.QueryRow(`SELECT id FROM roles where name=$1 and is_global=true`,roleName).Scan(&roleId)
		if err!=nil{
			log.Printf(" Could not find role: %s. Skipping...\n", roleName)
			continue
		}

		for _,permissionName:=range permissionNames{
			var permId int64
			err:=conn.QueryRow(`SELECT id FROM permissions WHERE name=$1`,permissionName).Scan(&permId)
			if err!=nil{
				log.Printf("Could not find permission: %s. Skipping...\n", permissionName)
				continue
			}

			_,err=conn.Exec(`INSERT INTO role_permissions (role_id,permission_id)
			VALUES($1,$2)
			ON CONFLICT (role_id,permission_id) DO NOTHING;
			`,roleId,permId)

			if err!=nil{
				log.Printf(" Error inserting role-permission (%s - %s): %v\n", roleName, permissionName, err)
			}
		}
	}
	return nil
}