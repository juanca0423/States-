// Package db
package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DBConnection() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			PrepareStmt:            false,
			SkipDefaultTransaction: true,
		})

		if err == nil {
			// CONEXIÓN EXITOSA
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetMaxOpenConns(20)
			sqlDB.SetConnMaxLifetime(time.Minute * 10)

			DB = db
			fmt.Println("🚀 Conexión optimizada con PgBouncer exitosa")
			return
		}

		// Si llegamos aquí, es que hubo un error
		fmt.Printf("⚠️ Intento %d/5 fallido: %v. Reintentando en 2s...\n", i, err)
		time.Sleep(3 * time.Second)
	}

	// Si después de 5 intentos no conectó, cerramos la app
	log.Fatalf("❌ No se pudo conectar a la DB después de 5 intentos: %v", err)
}
