// Package rutas
package rutas

import (
	"ef/ctrl"
	"ef/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpRutas(app *fiber.App) {
	// PÚBLICAS
	app.Get("/", ctrl.GetIndex)

	app.Get("/loguin", ctrl.GetLoguin)
	app.Post("/loguin", ctrl.PostHandleLogin)

	app.Get("/register", ctrl.GetRejistro)
	app.Post("/register", ctrl.RegisterHandler)

	app.Get("/about", ctrl.GetAbout)

	// En ef/rutas/rutas.go
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("jwt")
		return c.Redirect("/")
	})

	// Ruta para herramienta comercial
	app.Get("/eeff", middleware.AuthRequired, ctrl.HojaTrabajo)
	app.Post("/estados", middleware.AuthRequired, ctrl.GenEstados)
	app.Post("/exportar-excel-todo", middleware.AuthRequired, ctrl.ExportarExcelTodo)
	app.Post("/api/pagos/qpay-webhook", ctrl.WebhookQPayPro)
	app.Get("/perfil", middleware.AuthRequired, ctrl.GetPerfil)
	// PROTEGIDAS (Requieren estar logueado)
	pago := app.Group("/pago", middleware.AuthRequired)
	pago.Get("/exito", func(c *fiber.Ctx) error {
		return c.Render("exito", fiber.Map{"Title": "¡Suscripción Activa!"})
	})

	// Ruta para la nueva herramienta industrial
	app.Post("/costos", middleware.AuthRequired, ctrl.GenCostoProduccion)
	app.Get("/costosform", middleware.AuthRequired, ctrl.HojaTrabajocosto)

	app.Get("/soport", middleware.AuthRequired, ctrl.GetSoporte)
	app.Post("/soporte/enviar", middleware.AuthRequired, ctrl.PostConsulta)

	// API Y DASHBOARDS (Roles específicos)
	api := app.Group("/api", middleware.AuthRequired)
	admin := api.Group("/admin", middleware.Require("admin"))
	admin.Get("/soporte", ctrl.GetAdminSoporte)
	admin.Post("/soporte/responder/:id", ctrl.PostResponder)
	admin.Get("/dashboard", ctrl.GetAdmDash)
	admin.Get("/usuario/:id", ctrl.GetUsuarioDetalle)
	// ... dentro de SetUpRutas ...
	admin.Get("/crearcuenta", ctrl.VerPanelCuentas)
	admin.Post("/crearcuenta", ctrl.CrearCuenta)
	admin.Get("/eliminar-cuenta/:codigo", ctrl.EliminarCuenta)
	admin.Post("/editar-cuenta", ctrl.PostEditarCuenta)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("error-404", fiber.Map{
			"Titulo": "Página no encontrada",
		})
	})
}
