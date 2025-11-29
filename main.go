package main

import (
	"fmt"
	"models"
	"db"
	"github.com/gofiber/template/handlebars/v3"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// database connection
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})

	engine := handlebars.New("./views", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/layout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		},
		"layouts/main")
	})
	app.Listen(":3000")
	fmt.Println("Serbidor corriemdo en el puerto: 3000")
}

