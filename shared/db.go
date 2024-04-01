package shared

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
		log.Printf("failed to insert %v", result.Error)
		return result.Error
	}

	log.Printf("inserted successfully")
	return nil
}

// func Query(query interface{}, result interface{}) (interface{}, error) {
// 	db := Connect()

// 	if err := db.Where(query).Find(result).Error; err != nil {
// 		log.Printf("failed to query %v", err)
// 		return nil, err
// 	}
// 	return nil, nil
// }

// func Update(record interface{}, updateMap map[string]interface{}) error {
// 	db := Connect()

// 	if err := db.Model(record).Where("...").Updates(updateMap).Error; err != nil {
// 		log.Printf("failed to update %v", err)
// 		return err
// 	}

// 	log.Printf("update successful")
// 	return nil
// }

// func Delete(record interface{}) error {
// 	db := Connect()

// 	if err := db.Where("...").Delete(record).Error; err != nil {
// 		log.Printf("failed to delete %v", err)
// 		return err
// 	}

// 	log.Printf("deletion successful")
// 	return nil
// }
