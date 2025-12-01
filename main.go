package main

import (
	"fmt"

	"ef/db"
	"ef/models"
	"ef/rutas"
)

func main() {
	// database connection
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})
	app := rutas.InitRouter()
	app.Listen(":3000")
	fmt.Println("Serbidor corriemdo en el puerto: 3000")
}
