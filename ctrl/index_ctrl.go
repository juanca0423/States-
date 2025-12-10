package ctrl

import (
	"github.com/gofiber/fiber/v2"
)

func GetIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Index",
	})
}

func GetLoguin(c *fiber.Ctx) error {
	return c.Render("loguin", fiber.Map{
		"Title": "Loguin",
	})
}

func GetAbout(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Que Hacemos",
	})
}

func GetRejistro(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "registro de usuarios",
	})
}

func PostRegistro(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "usuario registrado",
	})
}

func GetLayout(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}
