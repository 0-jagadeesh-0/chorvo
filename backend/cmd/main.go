package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0-jagadeesh-0/chorvo/database"

	"github.com/0-jagadeesh-0/chorvo/router"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}

// @title Chorvo API
// @version 1.0
// @description API documentation for Chorvo project management tool
// @host localhost:8080
// @BasePath /api/v1
func main() {

	database.ConnectDB()

	appRouter := router.SetupRouter()

	port := os.Getenv("PORT")

	host:= os.Getenv("HOST")
	
	log.Printf("Server running on http://%s%s", host, port)
	if err := http.ListenAndServe(port, appRouter); err != nil {
		log.Fatal(err)
	}
  } 