package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Nombre    string    `gorm:"not null" json:"nombre"`
	Apellido  string    `gorm:"not null" json:"apellido"`
	Email     string    `gorm:"not null; uniqueIndex" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	FechaPago time.Time `gorm: json:"fechapago"`
	Role      string    `gorm:"not null" json:"role"`
}
