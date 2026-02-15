// Package ctrl elcontrolador para generar los estados financieros
package ctrl

import (
	"strconv"

	"ef/config"
	"ef/help"
	"ef/models"

	"github.com/gofiber/fiber/v2"
)

type KV struct {
	Key   string
	Value models.HtString
}

func GenEstados(c *fiber.Ctx) error {
	form, err := c.Request().MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "multipart error"})
	}

	Balance := make(map[string]float64)

	for k, vv := range form.Value {
		if len(vv) > 0 {
			if f, err := strconv.ParseFloat(vv[0], 64); err == nil {
				if f != 0 {
					Balance[k] = f
				}
			}
		}
	}

	Balan := config.DividirCuentas(Balance)

	DBHojadetravajo := help.HojaDeTrabajo(Balance)
	DBResultados, TotResu := help.Resultados(Balance, 0) // 0 indica que es comercial
	DBBalnce, TotBalance := help.GenerarTodoElBalance(Balan, TotResu.UtilidadNeta)

	dashboard := help.GenerarDashboard(
		TotBalance.ActivoCorriente,
		TotBalance.PasivoCorriente,
		TotBalance.Inventario,
		TotResu.CostoVentas,
		TotBalance.ActivoTotal,
		TotBalance.PasivoTotal,
	)

	// 5. Renderizado
	return c.Render("estados", fiber.Map{
		"Title":       "Estados Financieros",
		"keys":        DBHojadetravajo,
		"resultados":  DBResultados,
		"BalanceRows": DBBalnce,
		"indices":     dashboard,
	})
}

func GenCostoProduccion(c *fiber.Ctx) error {
	// 1. Capturar datos del formulario industrial
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error en el formulario")
	}

	// 2. Convertir y Limpiar datos
	balanceRaw := make(map[string]float64)
	for k, vv := range form.Value {
		if len(vv) > 0 {
			if f, err := strconv.ParseFloat(vv[0], 64); err == nil {
				balanceRaw[k] = f
			}
		}
	}
	balan := config.DividirCuentas(balanceRaw)

	// 3. Generar la Hoja de Trabajo de 9 Columnas (Llamando a tu nueva función en help)
	hoja9Cols := help.HojaDeTrabajoIndustrial(balanceRaw)

	// 4. Cálculos para el Reporte de Costos (Usa los códigos de tu nomenclatura)
	getV := func(cod string) float64 { return balanceRaw[cod] }

	invInicialMP := getV("311001")
	comprasMP := getV("311002")
	gastosMP := getV("311003")
	devMP := getV("311004")
	invFinalMP := getV("111303") // Viene del Activo

	mod := getV("312001") // Mano de Obra Directa

	// Gastos Indirectos (CIF)
	cif := getV("312002") + getV("312003") + getV("313001") + getV("313002") + getV("313003")

	// Fórmulas Matemáticas
	comprasNetas := comprasMP + gastosMP - devMP
	mpConsumida := (invInicialMP + comprasNetas) - invFinalMP
	costoPrimo := mpConsumida + mod
	costoProduccion := costoPrimo + cif
	// Primero calculas el costo de producción con las fórmulas que vimos
	vCostoProd := mpConsumida + mod + cif

	// Se lo pasas a la función para que arme el Estado de Resultados Industrial
	DBResultados, TotResu := help.Resultados(balanceRaw, vCostoProd)
	DBBalnce, _ := help.GenerarTodoElBalance(balan, TotResu.UtilidadNeta)

	// 5. Empaquetar resultados para el HBS
	resCalculados := map[string]string{
		"InvInicialMP":    config.FCont(invInicialMP),
		"ComprasNetasMP":  config.FCont(comprasNetas),
		"MPDisponible":    config.FCont(invInicialMP + comprasNetas),
		"InvFinalMP":      config.FCont(invFinalMP),
		"MPConsumida":     config.FCont(mpConsumida),
		"MOD":             config.FCont(mod),
		"CostoPrimo":      config.FCont(costoPrimo),
		"CIF":             config.FCont(cif),
		"CostoProduccion": config.FCont(costoProduccion),
	}

	return c.Render("costos", fiber.Map{
		"Title":      "Auditoría Industrial - Costos",
		"filasHoja":  hoja9Cols,
		"Resultados": DBResultados,
		"balanceRaw": DBBalnce,
		"res":        resCalculados,
	})
}
