package db

// staff and roles and permissions

func CreateRolesTable() error {
	query := `CREATE TABLE IF NOT EXISTS roles(
	id SERIAL PRIMARY KEY,
    code VARCHAR(20),
	name VARCHAR(100) NOT NULL,
	description TEXT,
	restaurant_id INTEGER,
	is_global BOOLEAN DEFAULT FALSE,
	level INTEGER DEFAULT 0,
	is_custom BOOLEAN DEFAULT TRUE,
	manager_id INTEGER,  -- if null that means by restaurant owner
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
	FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
	UNIQUE(name,restaurant_id)
	)`

	_, err := DB.Exec(query)
	return err
}

// permissions like general permissions: view_order,see_reports and authority permissions like manage, delete etc;
func CreatePermissionsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
	`
	_, err := DB.Exec(query)
	return err
}

// This handles general permissions like view_orders, edit_menu, etc.
func CreateRolePermissionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE(role_id, permission_id)
   )`
	_, err := DB.Exec(query)
	return err
}

// who(roles) have authority permission on whome(roles)
func CreateRolesTargetPermission() error {
	query := `CREATE TABLE IF NOT EXISTS role_target_permissions (
    id SERIAL PRIMARY KEY,
    actor_role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    action_permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    target_role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE(actor_role_id, action_permission_id, target_role_id)
)`
	_, err := DB.Exec(query)
	return err
}

func CreateStaffRoles() error {
	query := `CREATE TABLE IF NOT EXISTS staff_roles (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    is_approved INTEGER DEFAULT 0,  -- 0 for that he is requested, -1 for rejected, 1 for accepted, -2 for blocked
    manager_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(user_id, restaurant_id,role_id)
   )`
	_, err := DB.Exec(query)
	return err
}
