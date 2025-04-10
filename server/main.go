package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/harshgupta9473/restaurantmanagement/db"
	"github.com/harshgupta9473/restaurantmanagement/routes"
	"github.com/harshgupta9473/restaurantmanagement/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	log.Println("ENV loaded")

	// Initialize database
	err = db.InitDB()
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	log.Println("Database initialized")

	// Create tables
	err = db.CreateAllTable()
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Tables created")

	// Load secrets
	utils.LoadSecrets()

	// Setup routes
	router := routes.SetupRoutes()

	// Configure HTTP server
	s := &http.Server{
		Addr:         ":3001",
		Handler:      &router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Listening on port 3001")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	log.Println("Received signal to terminate:", sig)

	// Timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("Server exited properly")
}
