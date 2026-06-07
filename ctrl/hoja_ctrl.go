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
		"Title":   "Ingreso de Datos Empresa Comercial",
		"Cuentas": string(cuentasJSON),
		"EsCosto": esCosto,
	})
}

func HojaTrabajoCosto(c *fiber.Ctx) error {
	esCosto := c.QueryBool("costos", true)
	cuentasData := config.ObtenerCuentas(esCosto)

	cuentasJSON, err := json.Marshal(cuentasData)
	if err != nil {
		return c.Status(500).SendString("Error serializando nomenclatura")
	}

	return c.Render("costosform", fiber.Map{
		"Title":   "Ingreso de Datos Empresa Industrial",
		"Cuentas": string(cuentasJSON),
		"EsCosto": esCosto,
	})
}
