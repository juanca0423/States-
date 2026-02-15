package main

import (
	"log"

	"ef/config"
	"ef/db"
	"ef/middleware"
	"ef/rutas"
	"ef/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/handlebars/v3"
)

func main() {
	// 1. Intentar conectar a la DB con los reintentos que configuramos
	db.DBConnection()
	// 2. Verificar que la instancia de DB no sea nil antes de usarla
	if db.DB == nil {
		log.Fatal("❌ La instancia de base de datos es nula. Abortando...")
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatal("❌ Error al obtener el objeto sql.DB de GORM:", err)
	}
	// 3. Cargar nomenclatura (Garantiza que tengamos las 66 cuentas en memoria)
	err = config.CargarNomenclaturaDesdeDB(sqlDB)
	if err != nil {
		// Logueamos el error pero no matamos la app, para poder entrar al admin y arreglarlo
		log.Printf("⚠️ No se pudo cargar la nomenclatura: %v", err)
	}
	// 4. Configuración de Engine de Vistas
	engine := handlebars.New("./views", ".hbs")
	utils.RegistrarHelpers(engine)
	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		// ProxyHeader: "X-Forwaded-For",
	})
	// 5. Middlewares y Rutas
	app.Static("/static", "./estaticos")
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		// Optimizamos la comprobación del path estático
		if len(c.Path()) >= 7 && c.Path()[:7] == "/static" {
			return c.Next()
		}
		return middleware.AuthStatus(c)
	})
	rutas.SetUpRutas(app)
	// 6. Encendido del servidor
	log.Fatal(app.Listen(":3000"))
}
