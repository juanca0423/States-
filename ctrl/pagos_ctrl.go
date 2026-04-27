package ctrl

import (
	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

func WebhookQPayPro(c *fiber.Ctx) error {
	// 1. QPayPro envía los datos en el Body
	type QPayResponse struct {
		UserID     uint    `json:"user_id"`
		Monto      float64 `json:"amount"`
		Estado     string  `json:"status"`
		Referencia string  `json:"transaction_id"`
	}

	var data QPayResponse
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Error en datos")
	}

	// 2. Si el pago fue exitoso, actualizamos al usuario
	if data.Estado == "SUCCESS" {
		// Actualizar suscripción del usuario
		db.DB.Model(&models.User{}).Where("id = ?", data.UserID).
			Update("suscripcion_activa", true)

		// Registrar la transacción
		nuevaTrans := models.Transaccion{
			UserID:     data.UserID,
			Monto:      data.Monto,
			Estado:     data.Estado,
			Referencia: data.Referencia,
			Pasarela:   "qpaypro",
		}
		db.DB.Create(&nuevaTrans)
	}

	return c.SendStatus(200) // Siempre responder 200 a la pasarela
}
