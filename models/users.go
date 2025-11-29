package models

import "gorm.io/gorm"
type User struct{
	gorm.Model

	First_ Name string `gorm:"not null" json:"first_name"`
	Last_Name string `gorm:"not null" json:"last_name"`
	Email string `gorm:"not null; uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role string `gorm:"not null" json:"role"`
}

