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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ENV LOADED")
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db is initialiesd")
	err = db.CreateAllTable()
	if err != nil {
		log.Fatal("Error creating tables %w", err)
	}
	log.Println("Tables created")

	utils.LoadSecrets()

	router := routes.SetupRoutes()

	s := &http.Server{
		Addr:         ":3001",
		Handler:      &router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)

}
