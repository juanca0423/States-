// Package ctrl
package ctrl

import (
	"time"

	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

func GetAdmDash(c *fiber.Ctx) error {
	var totalUsuarios int64
	var usuarios []models.User
	var consultasPendientes int64

	// Contar usuarios totales
	db.DB.Model(&models.User{}).Count(&totalUsuarios)

	// Traer los últimos 10 registrados
	db.DB.Order("created_at desc").Limit(10).Find(&usuarios)

	// Contar cuántas dudas hay sin responder
	db.DB.Model(&models.Mensaje{}).Where("estado = ?", "Pendiente").Count(&consultasPendientes)

	return c.Render("admin_dashboard", fiber.Map{
		"Title":               "Panel de Control Administrativo",
		"Total":               totalUsuarios,
		"Usuarios":            usuarios,
		"ConsultasPendientes": consultasPendientes,
	})
}

func responderConsulta(id, respuesta string) error {
	return db.DB.Model(&models.Mensaje{}).Where("id = ?", id).Updates(map[string]any{
		"respuesta": respuesta,
		"estado":    "Resuelto",
	}).Error
}

// GetAdminSoporte muestra todas las dudas de todos los usuarios
func GetAdminSoporte(c *fiber.Ctx) error {
	var consultas []models.Mensaje
	// Traemos las consultas incluyendo los datos del Usuario (Nombre/Email)
	db.DB.Preload("User").Order("created_at desc").Find(&consultas)

	return c.Render("admin_soporte", fiber.Map{
		"Title":     "Panel de Expertos - Consultas",
		"Consultas": consultas,
	})
}

// ResponderConsulta guarda la respuesta y marca como resuelto
func ResponderConsulta(c *fiber.Ctx) error {
	id := c.Params("id")
	respuesta := c.FormValue("respuesta")

	if err := responderConsulta(id, respuesta); err != nil {
		return c.Status(500).SendString("No se pudo responder")
	}
	return c.Redirect("/api/admin/soporte")
}

// PostResponder guarda la respuesta y marca como resuelto
func PostResponder(c *fiber.Ctx) error {
	id := c.Params("id")
	respuesta := c.FormValue("respuesta")

	if respuesta == "" {
		return c.Redirect("/api/admin/soporte?error=vacio")
	}

	if err := responderConsulta(id, respuesta); err != nil {
		return c.Status(500).SendString("No se pudo responder")
	}
	return c.Redirect("/api/admin/soporte?exito=respondido")
}

func PostActivarUsuario(c *fiber.Ctx) error {
	id := c.Params("id")

	// Actualizamos el estado y extendemos la fecha de fin (opcional)
	err := db.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]any{
		"suscripcion_activa": true,
		"fecha_fin_prueba":   time.Now().AddDate(1, 0, 0), // Le damos 1 año de acceso
	}).Error
	if err != nil {
		return c.Status(500).SendString("No se pudo actualizar")
	}

	return c.SendStatus(200)
}

func GetUsuarioDetalle(c *fiber.Ctx) error {
	id := c.Params("id")
	var u models.User

	if err := db.DB.First(&u, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No encontrado"})
	}

	// Retornamos un JSON estructurado explícitamente con las claves exactas en mayúsculas que espera miscript.js
	return c.Status(200).JSON(fiber.Map{
		"ID":                u.ID,
		"Nombre":            u.Nombre,
		"Apellido":          u.Apellido,
		"Email":             u.Email,
		"Role":              u.Role,
		"FechaFinPrueba":    u.FechaFinPrueba,
		"SuscripcionActiva": u.SuscripcionActiva,
	})
}
