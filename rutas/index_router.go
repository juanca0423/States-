package rutas

import (
	"ef/ctrl"

	"github.com/gofiber/fiber/v2"
)

func SetUpRutas(app *fiber.App) {
	app.Get("/", ctrl.GetIndex)
	app.Get("/loguin", ctrl.GetLoguin)
	app.Get("/register", ctrl.GetRejistro)
	app.Get("/about", ctrl.GetAbout)

	app.Post("/register", ctrl.PostRegistro)

	app.Get("/layout", ctrl.GetLayout)
}
