package main

import (
	"fmt"

	"ef/db"
	"ef/models"
	"ef/rutas"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/handlebars/v3"
)

func main() {
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})

	engine := handlebars.New("./views", ".hbs")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./estaticos")
	app.Use(logger.New())
	rutas.SetUpRutas(app)

	app.Listen(":3000")
	fmt.Println("Serbidor corriemdo en el puerto: 3000")
}
