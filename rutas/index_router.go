package rutas

import (
	"github.com/gofiber/fiber/v2"
)

func GetIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Index",
	})
}

func GetLoguin(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Loguin",
	})
}

func GetRejistro(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "registro de usuarios",
	})
}

func PostRegistro(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "usuario registrado",
	})
}
