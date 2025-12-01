package models

import (
	"gorm.io/gorm"
	"ef/roles"
)

type User struct{
	gorm.Model

	First_Name string `gorm:"not null" json:"first_name"`
	Last_Name string `gorm:"not null" json:"last_name"`
	Email string `gorm:"not null; uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role roles.RolUser `gorm:"not null" json:"role"`
}

