package db

import (
	"fmt"
	"os"

	"github.com/Kalmera74/Shorty/internal/shortener"
	"github.com/Kalmera74/Shorty/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AutoMigrate(dbConn *gorm.DB) error {
	models := []interface{}{
		&user.UserModel{},
		&shortener.ShortModel{},
	}

	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to auto-migrate model: %v", err)
		}
	}
	return nil
}
