package db

import (
	"errors"
	"log"
	"time"

	"ef/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin() {
	// 1. Definimos el correo del admin en una variable para no fallar
	adminEmail := "juanca0423@gmail.com"
	var admin models.User
	// 2. Buscamos por el correo real que vamos a usar
	err := DB.Where("email = ?", adminEmail).First(&admin).Error
	// 3. Si el error es "Record Not Found", lo creamos
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		nuevoAdmin := models.User{
			Nombre:            "Juan Carlos",
			Apellido:          "Perez Castro",
			Email:             adminEmail, // Usamos la misma variable
			Pase:              string(hash),
			Role:              "admin",
			SuscripcionActiva: true,
			FechaFinPrueba:    time.Now().AddDate(99, 0, 0),
		}
		if errCreate := DB.Create(&nuevoAdmin).Error; errCreate != nil {
			log.Printf("No se pudo crear el admin: %v", errCreate)
		} else {
			log.Println("✅ Usuario administrador creado por defecto con éxito")
		}
	} else if err != nil {
		log.Printf("Error al consultar la BD: %v", err)
	} else {
		log.Println("ℹ️ El administrador ya existe en el sistema")
	}
}
