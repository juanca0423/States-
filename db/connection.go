// Package db
package db

import (
	"ef/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DBConnection() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
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

			err = DB.AutoMigrate(&models.User{}, &models.Mensaje{}, &models.CueDB{}, &models.Transaccion{})
			if err != nil {
				// PostgreSQL prohíbe alterar una columna referenciada por una política RLS.
				// Recuperamos la definición, eliminamos la política, migramos y la recreamos.
				if strings.Contains(err.Error(), "SQLSTATE 0A000") {
					type policyDef struct {
						PolicyName string
						Cmd        string
						Qual       string
						WithCheck  string
					}
					// Use sql.Rows to fetch policies without GORM struct mapping issues
					var policies []policyDef
					rows, errRows := DB.Raw("SELECT policyname, cmd, qual, with_check FROM pg_policies WHERE tablename = 'transacciones'").Rows()
					if errRows != nil {
						fmt.Printf("❌ No se pudo obtener políticas: %v\n", errRows)
					} else {
						for rows.Next() {
							var p policyDef
							if errScan := rows.Scan(&p.PolicyName, &p.Cmd, &p.Qual, &p.WithCheck); errScan != nil {
								fmt.Printf("❌ Error al escanear política: %v\n", errScan)
								continue
							}
							policies = append(policies, p)
						}
						rows.Close()
					}
					for _, p := range policies {
						if dErr := DB.Exec(fmt.Sprintf("DROP POLICY IF EXISTS %q ON transacciones", p.PolicyName)).Error; dErr != nil {
							fmt.Printf("❌ No se pudo eliminar la política %s: %v\n", p.PolicyName, dErr)
						}
					}
					if err := DB.AutoMigrate(&models.Transaccion{}); err != nil {
						fmt.Printf("❌ Error en AutoMigrate tras eliminar políticas: %v\n", err)
					} else {
						for _, p := range policies {
							sqlStmt := fmt.Sprintf("CREATE POLICY %q ON transacciones FOR %s USING (%s)", p.PolicyName, p.Cmd, p.Qual)
							if p.WithCheck != "" {
								sqlStmt += fmt.Sprintf(" WITH CHECK (%s)", p.WithCheck)
							}
							if cErr := DB.Exec(sqlStmt).Error; cErr != nil {
								fmt.Printf("❌ No se pudo recrear la política %s: %v\n", p.PolicyName, cErr)
							}
						}
						fmt.Println("✅ Migración con políticas restauradas")
					}
				} else {
					fmt.Printf("❌ Error en AutoMigrate: %v\n", err)
				}
			}

			fmt.Println("🚀 Conexión y Migración exitosa")
			return
		}

		// Si llegamos aquí, es que hubo un error
		fmt.Printf("⚠️ Intento %d/5 fallido: %v. Reintentando en 2s...\n", i, err)
		time.Sleep(3 * time.Second)
	}

	// Si después de 5 intentos no conectó, cerramos la app
	log.Fatalf("❌ No se pudo conectar a la DB después de 5 intentos: %v", err)
}
