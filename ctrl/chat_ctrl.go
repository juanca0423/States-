// Package ctrl
package ctrl

import (
	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

func GetSoporte(c *fiber.Ctx) error {
	uidRaw := c.Locals("uid")
	if uidRaw == nil {
		return c.Redirect("/loguin")
	}

	// Conversión segura: intentamos pasarlo a uint, si falla redirigimos
	uid, ok := uidRaw.(uint)
	if !ok {
		// A veces los JWT guardan los números como float64
		if f, isFloat := uidRaw.(float64); isFloat {
			uid = uint(f)
		} else {
			return c.Redirect("/loguin")
		}
	}

	var consultas []models.Mensaje
	db.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&consultas)

	return c.Render("soport", fiber.Map{
		"Title":     "Soporte Técnico",
		"Consultas": consultas,
		"Enviado":   c.Query("enviado") == "true",
	})
}

func PostConsulta(c *fiber.Ctx) error {
	uidRaw := c.Locals("uid")
	if uidRaw == nil {
		return c.Redirect("/loguin")
	}

	uid, _ := uidRaw.(uint) // O la conversión float64 de arriba
	texto := c.FormValue("consulta")

	if texto == "" {
		return c.Redirect("/soport?error=vacio")
	}
	nuevoMensaje := models.Mensaje{
		UserID:   uid,
		Consulta: texto,
	}

	if err := db.DB.Create(&nuevoMensaje).Error; err != nil {
		return c.Status(500).SendString("Error al enviar consulta")
	}

	return c.Redirect("/soport?enviado=true")
}
