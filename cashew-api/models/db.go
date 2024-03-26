package models

import (
	"log"
	// "testing"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "host=localhost user=postgres password=polarBear$02 dbname=cashew_db port=5432 sslmode=disable TimeZone=America/New_York"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("connected to database successfully")
	return db
}

func Insert(record interface{}) error {
	db := Connect()
	result := db.Create(record)
	if result.Error != nil {
		log.Printf("failed to insert feedback %v", result.Error)
		return result.Error
	}

	log.Printf("feedback inserted successfully")
	return nil
}
