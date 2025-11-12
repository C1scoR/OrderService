package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	Item     string `gorm:"unique"`
	Quantity int32
}
type User struct {
	ID int `gorm:"primaryKey"` //надо донастроить, что нужный мне тип - это int
}

// BeforeCreate — хук, вызываемый GORM до создания записи
func (m *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
