package models

import (
	"log"

	. "github.com/brianykl/cashew/shared"
)

func Query(query interface{}, result interface{}) (interface{}, error) {
	db := Connect()

	if err := db.Where(query).Find(result).Error; err != nil {
		log.Printf("failed to query %v", err)
		return nil, err
	}
	return nil, nil
}

func Update(record interface{}, updateMap map[string]interface{}) error {
	db := Connect()

	if err := db.Model(record).Where("...").Updates(updateMap).Error; err != nil {
		log.Printf("failed to update %v", err)
		return err
	}

	log.Printf("update successful")
	return nil
}

func Delete(record interface{}) error {
	db := Connect()

	if err := db.Where("...").Delete(record).Error; err != nil {
		log.Printf("failed to delete %v", err)
		return err
	}

	log.Printf("deletion successful")
	return nil
}
