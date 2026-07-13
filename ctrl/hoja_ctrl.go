package ctrl

import (
	"encoding/json"

	"ef/config"

	"github.com/gofiber/fiber/v2"
)

func renderHoja(c *fiber.Ctx, esCosto bool, plantilla string) error {
	cuentasData := config.ObtenerCuentas(esCosto)
	cuentasJSON, err := json.Marshal(cuentasData)
	if err != nil {
		return c.Status(500).SendString("Error serializando nomenclatura")
	}
	titulo := "Ingreso de Datos Empresa Comercial"
	if esCosto {
		titulo = "Ingreso de Datos Empresa Industrial"
	}
	return c.Render(plantilla, fiber.Map{
		"Title":   titulo,
		"Cuentas": string(cuentasJSON),
		"EsCosto": esCosto,
	})
}

func HojaTrabajo(c *fiber.Ctx) error {
	return renderHoja(c, c.QueryBool("costos", false), "eeffform")
}

func HojaTrabajoCosto(c *fiber.Ctx) error {
	return renderHoja(c, c.QueryBool("costos", true), "costosform")
}
