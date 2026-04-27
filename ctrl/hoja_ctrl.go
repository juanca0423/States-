package ctrl

import (
	"ef/config"

	"github.com/gofiber/fiber/v2"
)

func HojaTrabajo(c *fiber.Ctx) error {
	esCosto := c.QueryBool("costos", false)
	return c.Render("eeffform", fiber.Map{
		"Title":   "Ingreso de Datos",
		"Cuentas": config.ObtenerCuentas(esCosto),
		"EsCosto": esCosto,
	})
}

func HojaTrabajocosto(c *fiber.Ctx) error {
	esCosto := c.QueryBool("costos", true)
	return c.Render("costosform", fiber.Map{
		"Title":   "Ingreso de Datos",
		"Cuentas": config.ObtenerCuentas(esCosto),
		"EsCosto": esCosto,
	})
}
