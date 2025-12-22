package rutas

import (
	"ef/ctrl"
	"ef/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpRutas(app *fiber.App) {
	app.Get("/", ctrl.GetIndex)
	app.Get("/loguin", ctrl.GetLoguin)
	app.Get("/register", ctrl.GetRejistro)
	app.Get("/about", ctrl.GetAbout)

	app.Post("/register", ctrl.RegisterHandler)
	app.Post("loguin", ctrl.PostHandleLogin)
	app.Get("/layout", ctrl.GetLayout)
	app.Get("/eeff", ctrl.HojaTrabajo)
	app.Post("/estados", ctrl.GenEstados)

	api := app.Group("/api")
	// cualquier usuario autenticado
	user := api.Group("/user", middleware.Any("usuario", "cliente", "admin"))
	user.Get("/perfil", ctrl.GetPerfil)

	// cliente o admin
	cliente := api.Group("/cliente", middleware.Any("cliente", "admin"))
	cliente.Get("/pedidos", ctrl.PostRegistro)

	// solo administradores
	admin := api.Group("/admin", middleware.Require("admin"))
	admin.Get("/dashboard", ctrl.GetAdmDash)

}
