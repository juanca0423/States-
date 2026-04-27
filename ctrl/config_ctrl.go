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
	// 1. Capturamos el valor del select ("true" o "false")
	escostoStr := c.FormValue("escosto")
	// 2. Convertimos el string a booleano
	// ParseBool acepta: "1", "t", "T", "true", "TRUE", "True", "0", "f", "F", "false", "FALSE", "False"
	escosto, _ := strconv.ParseBool(escostoStr)
	codigo, _ := strconv.Atoi(codigoStr)
	var existe models.CueDB
	if db.DB.Where("codigo = ?", codigo).First(&existe).RowsAffected > 0 {
		return c.Redirect("/api/admin/crearcuenta?error=El código " + codigoStr + " ya existe.")
	}
	// 3. Asignamos la variable booleana ya convertida
	nueva := models.CueDB{
		Codigo:    codigo,
		Nombre:    nombre,
		Categoria: categoria,
		Saldo:     saldo,
		EsCosto:   escosto,
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
	codigoStr := c.FormValue("Codigo") // Viene del input hidden id="edit_codigo"

	var cuenta models.CueDB
	// Buscamos la cuenta. Importante: Convertir a int si es necesario o dejar que GORM lo haga
	if err := db.DB.First(&cuenta, "codigo = ?", codigoStr).Error; err != nil {
		fmt.Println("Error: No se encontró la cuenta con código:", codigoStr)
		return c.Redirect("/api/admin/crearcuenta?error=no_encontrada")
	}

	// Actualizamos los valores
	cuenta.Nombre = c.FormValue("nombre")
	cuenta.Categoria = c.FormValue("categoria")
	cuenta.Saldo = c.FormValue("saldo") // Debe coincidir con el name="Saldo" del modal

	escostoStr := c.FormValue("escosto")
	cuenta.EsCosto, _ = strconv.ParseBool(escostoStr)

	// Guardamos los cambios físicamente en la DB
	if err := db.DB.Save(&cuenta).Error; err != nil {
		return c.Redirect("/api/admin/crearcuenta?error=db_error")
	}

	// RECARGA DE MEMORIA: Esto es lo que hace que los 65 pasen a 66
	sqlDB, _ := db.DB.DB()
	config.CargarNomenclaturaDesdeDB(sqlDB)

	return c.Redirect("/api/admin/crearcuenta?success=actualizada")
}
