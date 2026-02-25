package models

type User struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Email       string
	PhoneNumber string
	Addresses   []Address
}
