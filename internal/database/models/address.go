package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
	UserID  string
}
