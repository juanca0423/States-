package main

import (
	"fmt"

	"ef/db"
	"ef/models"
	"ef/rutas"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v3"
)

func main() {
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})
	engine := handlebars.New("./views", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", rutas.GetIndex)
	app.Get("/loguin", rutas.GetLoguin)

	app.Post("/register", rutas.PostRegistro)
	app.Listen(":3000")
	fmt.Println("Serbidor corriemdo en el puerto: 3000")
}
