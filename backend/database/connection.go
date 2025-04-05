package database

import (
	"fmt"
	"log"

	"github.com/0-jagadeesh-0/chorvo/config"
	"github.com/0-jagadeesh-0/chorvo/internal/models/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// 🧱 Create the users table if it doesn't exist
    err = DB.AutoMigrate(&entities.User{})
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
}
