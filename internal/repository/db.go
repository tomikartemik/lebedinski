package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lebedinski/internal/model"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		model.Category{},
		model.Item{},
		model.Photo{},
		model.Size{},
		model.Cart{},
		model.Order{},
		model.CartItem{},
		model.Top{},
		model.PromoCode{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
