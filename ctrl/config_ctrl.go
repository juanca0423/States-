// Package ctrl
package ctrl

import (
	"fmt"
	"strconv"

	"ef/config"
	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 1. Renderizar la página de gestión

// VerPanelCuentas Ver panel de cuentas
func VerPanelCuentas(c *fiber.Ctx) error {
	cuentas := config.ItemsCostos
	fmt.Printf("Mostrando panel: %d cuentas cargadas desde memoria\n", len(cuentas))
	return c.Render("crearcuenta", fiber.Map{
		"todasLasCuentas": cuentas,
		"error":           c.Query("error"),
		"success":         c.Query("success"),
		"UserName":        c.Locals("UserName"),
		"IsLogged":        true,
		"role":            "admin",
	})
}

// EliminarCuenta Eliminar cuenta
func EliminarCuenta(c *fiber.Ctx) error {
	// Obtenemos el código desde los parámetros de la URL
	codigo := c.Params("codigo")
	// 1. Eliminar de Supabase
	if err := db.DB.Where("codigo = ?", codigo).Delete(&models.CueDB{}).Error; err != nil {
		return c.Redirect("/api/admin/crearcuenta?error=No se pudo eliminar la cuenta")
	}
	// 2. RECARGAR MEMORIA (Fundamental para que desaparezca de los formularios)
	sqlDB, _ := db.DB.DB()
	config.CargarNomenclaturaDesdeDB(sqlDB)
	return c.Redirect("/api/admin/crearcuenta?success=Cuenta eliminada correctamente")
}

// CrearCuenta Crear cuenta

func CrearCuenta(c *fiber.Ctx) error {
	codigoStr := c.FormValue("codigo")
	nombre := c.FormValue("nombre")
	categoria := c.FormValue("categoria")
	saldo := c.FormValue("saldo")

	// Capturamos los booleanos
	escostoStr := c.FormValue("escosto")
	esvariableStr := c.FormValue("es_variable") // Nombre que pusimos en el select del HTML
	esefectivoStr := c.FormValue("es_efectivo")

	escosto, _ := strconv.ParseBool(escostoStr)
	esvariable, _ := strconv.ParseBool(esvariableStr)
	esefectivo, _ := strconv.ParseBool(esefectivoStr)
	codigo, _ := strconv.Atoi(codigoStr)

	var existe models.CueDB
	if db.DB.Where("codigo = ?", codigo).First(&existe).RowsAffected > 0 {
		return c.Redirect("/api/admin/crearcuenta?error=El código " + codigoStr + " ya existe.")
	}

	nueva := models.CueDB{
		Codigo:     codigo,
		Nombre:     nombre,
		Categoria:  categoria,
		Saldo:      saldo,
		EsCosto:    escosto,
		EsVariable: esvariable, // <-- Nuevo
		EsEfectivo: esefectivo, // <-- Nuevo
	}

	if err := db.DB.Create(&nueva).Error; err != nil {
		return c.Redirect("/api/admin/crearcuenta?error=Error al guardar")
	}

	sqlDB, _ := db.DB.DB()
	config.CargarNomenclaturaDesdeDB(sqlDB)
	return c.Redirect("/api/admin/crearcuenta?success=Cuenta '" + nombre + "' agregada con éxito")
}

// ActualizarCuenta Actualisar Cuentas
func ActualizarCuenta(c *fiber.Ctx) error {
	codigo := c.Params("id") // El ID que viene en la URL
	var cuenta models.CueDB
	// 1. Buscamos la cuenta
	if err := db.DB.Session(&gorm.Session{PrepareStmt: false}).First(&cuenta, "codigo = ?", codigo).Error; err != nil {
		return c.Status(404).SendString("Cuenta no encontrada")
	}
	// 2. Parseamos los nuevos datos del formulario
	if err := c.BodyParser(&cuenta); err != nil {
		return c.Status(400).SendString("Datos inválidos")
	}
	// 3. Guardamos los cambios
	db.DB.Session(&gorm.Session{PrepareStmt: false}).Save(&cuenta)
	return c.Redirect("/api/admin/crearcuenta")
}

// PostEditarCuenta Editar cuentas
func PostEditarCuenta(c *fiber.Ctx) error {
	codigoStr := c.FormValue("Codigo")

	var cuenta models.CueDB
	if err := db.DB.First(&cuenta, "codigo = ?", codigoStr).Error; err != nil {
		return c.Redirect("/api/admin/crearcuenta?error=no_encontrada")
	}

	cuenta.Nombre = c.FormValue("nombre")
	cuenta.Categoria = c.FormValue("categoria")
	cuenta.Saldo = c.FormValue("saldo")

	// Actualizamos los booleanos
	escostoStr := c.FormValue("escosto")
	esvariableStr := c.FormValue("es_variable")
	esefectivoStr := c.FormValue("es_efectivo")

	cuenta.EsCosto, _ = strconv.ParseBool(escostoStr)
	cuenta.EsVariable, _ = strconv.ParseBool(esvariableStr)
	cuenta.EsEfectivo, _ = strconv.ParseBool(esefectivoStr)

	if err := db.DB.Save(&cuenta).Error; err != nil {
		return c.Redirect("/api/admin/crearcuenta?error=db_error")
	}

	// Recarga de memoria
	sqlDB, _ := db.DB.DB()
	config.CargarNomenclaturaDesdeDB(sqlDB)

	return c.Redirect("/api/admin/crearcuenta?success=actualizada")
}
