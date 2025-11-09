package migrations

import (
	"fmt"
	"orderService/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.Migrator().CreateTable(&models.Order{})
	if err != nil {
		return err
	}
	has := db.Migrator().HasTable(&models.Order{})
	if has == false {
		return fmt.Errorf("The table ORDERS was not created throughout migrations")
	}
	return nil

}
