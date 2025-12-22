package ctrl

import (
	"ef/config"

	"github.com/gofiber/fiber/v2"
)

func HojaTrabajo(c *fiber.Ctx) error {
	Datos := config.CreaCuentas()
	return c.Render("eeffform", fiber.Map{
		"Title": "mi lista",
		"datos": Datos,
	})
}
