package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Nombre   string    `gorm:"not null" json:"nombre"`
	Apellido string    `gorm:"not null" json:"apellido"`
	Email    string    `gorm:"not null; uniqueIndex" json:"email"`
	Pase     string    `gorm:"not null" json:"pase"`
	Pago     time.Time `json:"pago"`
	Role     string    `gorm:"not null" json:"role"`
}
