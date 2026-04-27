// Package ctrl perfiles de Usuario
package ctrl

import (
	"fmt"
	"time"

	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

func GetPerfil(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(uint)
	if !ok {
		return c.Redirect("/login")
	}
	var u models.User
	var transacciones []models.Transaccion // Cambiado de TransaccionPagalo a Transaccion

	db.DB.First(&u, uid)

	// Usamos el slice 'transacciones' que declaramos arriba
	if err := db.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&transacciones).Error; err != nil {
		fmt.Println("Error:", err)
	}

	dias := max(0, int(time.Until(u.FechaFinPrueba).Hours()/24))
	return c.Render("perfil", fiber.Map{
		"Usuario":          u, // Esto permite usar {{Usuario.Role}}
		"DiasRestantes":    dias,
		"FechaVencimiento": u.FechaFinPrueba.Format("02/01/2006"),
		"Historial":        transacciones,
		"TextoDias":        fmt.Sprintf("%d días", dias),
		"ClaseAlerta":      c.Locals("ClaseAlerta"),
	})
}
