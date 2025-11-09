package migrations

import (
	"log"
	"orderService/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.Migrator().CreateTable(&models.Order{})
	if err != nil {
		log.Printf("Something went wrong during the migrations: %v", err)
	}
	has := db.Migrator().HasTable(&models.Order{})
	if has == false && err == nil {
		log.Println("The table ORDERS was not created throughout migrations")
	}
	//Напишем тестовые миграции
	testData := []models.Order{
		{Item: "Laptop", Quantity: 1},
		{Item: "Smartphone", Quantity: 2},
	}
	for _, record := range testData {
		if err := db.Create(&record).Error; err != nil {
			log.Printf("Something went wrong, when creating a test data: %v", err)
		}
	}
}
