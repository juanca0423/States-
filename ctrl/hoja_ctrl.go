package ctrl

import (
	"encoding/json"

	"ef/config"

	"github.com/gofiber/fiber/v2"
)

func HojaTrabajo(c *fiber.Ctx) error {
	esCosto := c.QueryBool("costos", false)
	cuentasData := config.ObtenerCuentas(esCosto)

	cuentasJSON, err := json.Marshal(cuentasData)
	if err != nil {
		return c.Status(500).SendString("Error serializando nomenclatura")
	}

	return c.Render("eeffform", fiber.Map{
		"Title":   "Ingreso de Datos",
		"Cuentas": string(cuentasJSON),
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
