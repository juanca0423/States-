// Package ctrl
package ctrl

import (
	"time"

	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardInfo(c *fiber.Ctx) error {
	// 1. Obtenemos el ID que el middleware JWT guardó en Locals
	uid, ok := c.Locals("uid").(uint)
	if !ok {
		return c.Redirect("/loguin")
	}

	var u models.User
	if err := db.DB.First(&u, uid).Error; err != nil {
		return c.Redirect("/loguin")
	}

	// 2. Cálculo de días de prueba
	diasRestantes := max(0, int(time.Until(u.FechaFinPrueba).Hours()/24))

	// 3. Renderizamos la vista (puedes pasar estos datos a tu vista de estados)
	return c.Render("eeff", fiber.Map{
		"Usuario":       u,
		"DiasTrial":     diasRestantes,
		"Suscrito":      u.SuscripcionActiva,
		"MostrarAlerta": diasRestantes <= 5 && !u.SuscripcionActiva, // Alerta cuando falten 5 días
	})
}
