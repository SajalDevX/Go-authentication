package models

import "gorm.io/gorm"

type User struct {
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Role string
	gorm.Model
}
