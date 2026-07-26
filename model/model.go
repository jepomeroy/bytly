package model

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Bytly struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Redirect string `json:"redirect" gorm:"not null"`
	Bytly    string `json:"bytly" gorm:"unique,not null"`
	Clicked  uint64 `json:"clicked"`
	Random   bool   `json:"random"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbPort, dbName)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			// Create the database if it doesn't exist
			createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
			err = db.Exec(createDBSQL).Error
			if err != nil {
				log.Fatal("Failed to create database:", err)
			}
		} else {
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
	}

	err = db.AutoMigrate(&Bytly{})
	if err != nil {
		fmt.Println(err)
	}

	// Close the initial connection and reconnect to the newly created database
	_, err = db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}
}
