package models

type Order struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	Item     string `gorm:"unique"`
	Quantity int32
}
