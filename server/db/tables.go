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

	return nil
}

func CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(100) NOT NULL,
	middle_name VARCHAR(100),
	last_name VARCHAR(100),
	email VARCHAR(255) NOT NULL UNIQUE ,
	password VARCHAR(255) NOT NULL,
	verified BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	blocked  BOOLEAN DEFAULT FALSE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateSuperAdmin() error {
	query := `CREATE TABLE IF NOT EXISTS super_admins (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`
	_, err := DB.Exec(query)
	return err
}

func CreateOTPTable() error {
	query := `CREATE TABLE IF NOT EXISTS otp_table(
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL UNIQUE,
	token VARCHAR(30) NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	verified BOOLEAN DEFAULT FALSE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantReqTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurantsreq (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		street_address VARCHAR(255),
		locality VARCHAR(100),
		city VARCHAR(100),
		state VARCHAR(100),
		postal_code VARCHAR(20),
		country VARCHAR(100),

		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		food_type VARCHAR(20) CHECK (food_type IN ('veg', 'non-veg', 'both')) NOT NULL,
		contact_number VARCHAR(20),
		contact_email VARCHAR(255),
		image_url TEXT,
		gst_number VARCHAR(30) not null,
		pan_number VARCHAR(20) not null,

		status VARCHAR(20) CHECK (status IN ('pending','approve','block')) NOT NULL DEFAULT 'pending',
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantsTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurants (
		id SERIAL PRIMARY KEY,

		owner_id INTEGER NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		street_address VARCHAR(255),
		locality VARCHAR(100),
		city VARCHAR(100),
		state VARCHAR(100),
		postal_code VARCHAR(20),
		country VARCHAR(100),
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		food_type VARCHAR(20) CHECK (food_type IN ('veg', 'non-veg', 'both')) NOT NULL,
		contact_number VARCHAR(20),
		contact_email VARCHAR(255),
		image_url TEXT,
		gst_number VARCHAR(30),
		pan_number VARCHAR(20),
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantDetailsTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurant_details (
		id SERIAL PRIMARY KEY,
		restaurant_id INTEGER UNIQUE NOT NULL,
		avg_cost_for_two INTEGER,
		opening_time TIME,
		closing_time TIME,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantImagesTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurant_images (
		id SERIAL PRIMARY KEY,
		restaurant_id INTEGER NOT NULL,
		image_url TEXT NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantReviewsTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurant_reviews (
		id SERIAL PRIMARY KEY,
		restaurant_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		rating INTEGER CHECK (rating BETWEEN 1 AND 5),
		review TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}

func CreateRestaurantTagsTable() error {
	query := `CREATE TABLE IF NOT EXISTS restaurant_tags (
		id SERIAL PRIMARY KEY,
		restaurant_id INTEGER NOT NULL,
		tag VARCHAR(50) NOT NULL,
		FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	return err
}
