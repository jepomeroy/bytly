package model

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// invalidCatalogName is the Postgres SQLSTATE code returned when the target
// database does not exist yet (e.g. on first run against a fresh server).
const invalidCatalogName = "3D000"

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == invalidCatalogName {
			if createErr := createDatabase(dbHost, dbUser, dbPassword, dbPort, dbName); createErr != nil {
				log.Fatalf("Failed to create database %q: %v\n", dbName, createErr)
			}

			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("Error connecting to the database after creating it: %v\n", err)
			}
		} else {
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
	}

	err = db.AutoMigrate(&Bytly{})
	if err != nil {
		fmt.Println(err)
	}
}

// createDatabase connects to the default "postgres" maintenance database and
// issues a CREATE DATABASE for dbName, since Postgres has no "CREATE IF NOT
// EXISTS" support and dbName can't exist yet for us to connect to directly.
func createDatabase(host, user, password, port, dbName string) error {
	maintenanceDSN := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=postgres sslmode=disable", host, user, password, port)

	maintenanceDB, err := gorm.Open(postgres.Open(maintenanceDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := maintenanceDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	return maintenanceDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)).Error
}
