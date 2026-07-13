// Package ctrl elcontrolador para generar los estados financieros
package ctrl

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"ef/config"
	"ef/help"
	"ef/models"
	_ "image/jpeg" // Esto registra el de JPEG por si acaso
	_ "image/png"  // Esto registra el decodificador de PNG

	"github.com/gofiber/fiber/v2"
)

// extraerBalance limpia y convierte los valores del formulario multipart a un mapa de montos.
func extraerBalance(form *multipart.Form) map[string]float64 {
	balance := make(map[string]float64)
	for k, vv := range form.Value {
		if len(vv) > 0 && vv[0] != "" {
			valStr := strings.ReplaceAll(vv[0], ",", "")
			valStr = strings.TrimSpace(valStr)
			if f, err := strconv.ParseFloat(valStr, 64); err == nil && f != 0 {
				balance[k] = f
			}
		}
	}
	return balance
}

func GenEstados(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		cuentasJSON, _ := json.Marshal(config.ObtenerCuentas(false))
		return c.Render("eeffform", fiber.Map{
			"Title":   "Ingreso de Datos Empresa Comercial",
			"Error":   "No se recibieron datos válidos. Por favor completa el formulario correctamente.",
			"Cuentas": string(cuentasJSON),
			"EsCosto": false,
		})
	}

	Balance := extraerBalance(form)

	Balan := config.DividirCuentas(Balance)

	DBHojadetravajo := help.HojaDeTrabajo(Balance)

	DBResultados, TotResu := help.Resultados(Balance, 0) // 0 indica que es comercial
	DBBalance, TotBalance := help.GenerarTodoElBalance(Balan, TotResu.UtilidadNeta)
	// ... después de obtener TotResu y TotBalance
	d := models.DatosDashboard{
		PasivoCorriente:   TotBalance.PasivoCorriente,
		InventarioInicial: TotBalance.InventarioInicial,
		ActivoCorriente:   TotBalance.ActivoCorriente,
		InventarioFinal:   TotBalance.Inventario,
		CostoVentas:       TotResu.CostoVentas,
		ActivoTotal:       TotBalance.ActivoTotal,
		PasivoTotal:       TotBalance.PasivoTotal,
		Ventas:            TotResu.Ventas,
		UtilidadNeta:      TotResu.UtilidadNeta,
		GastosFijos:       TotResu.GastosFijos,
		GastosVariables:   TotResu.GastosVariables,
		GastosNoEfectivo:  TotResu.GastosNoEfectivo,
	}

	// Llamada limpia
	dashboard, indicesTotales := help.GenerarDashboard(d)

	// Convertimos a JSON string para que JS no se confunda
	keysJSON, _ := json.Marshal(DBHojadetravajo)
	resultadosJSON, _ := json.Marshal(DBResultados)
	balanceJSON, _ := json.Marshal(DBBalance)
	indicesJSON, _ := json.Marshal(dashboard)
	// Calcula el porcentaje real: (Margen en Q / Ventas) * 100
	porcentajeMargen := 0.0
	if TotResu.Ventas > 0 {
		porcentajeMargen = (indicesTotales.MargenContribucion / TotResu.Ventas) * 100
	}

	// 5. Renderizado
	return c.Render("estados", fiber.Map{
		"Title":              "Estados Financieros",
		"keys":               DBHojadetravajo,
		"resultados":         DBResultados,
		"BalanceRows":        DBBalance,
		"indices":            dashboard,
		"nombreEmpresa":      "JC Dev Systems",
		"keysjson":           string(keysJSON),
		"resultadosjson":     string(resultadosJSON),
		"BalanceRowsjson":    string(balanceJSON),
		"indicesjson":        string(indicesJSON),
		"VentasNetas":        TotResu.VentasNetas,
		"CostosFijos":        TotResu.GastosFijos,
		"CostosVar":          TotResu.GastosVariables,
		"CostosTotales":      indicesTotales.CostosTotales,
		"MargenContribucion": fmt.Sprintf("%.2f", porcentajeMargen),
		"PuntoE":             indicesTotales.PuntoEContable,
		"PuntoEfe":           indicesTotales.PuntoECaja,
		"GastosNoEfectivo":   TotResu.GastosNoEfectivo,
	})
}

func GenCostoProduccion(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		cuentasJSON, _ := json.Marshal(config.ObtenerCuentas(true))
		return c.Render("costosform", fiber.Map{
			"Title":   "Ingreso de Datos Empresa Industrial",
			"Error":   "No se recibieron datos válidos. Por favor completa el formulario correctamente.",
			"Cuentas": string(cuentasJSON),
			"EsCosto": true,
		})
	}

	// 2. Convertir y Limpiar datos
	balanceRaw := extraerBalance(form)
	balan := config.DividirCuentas(balanceRaw)

	// 3. Generar la Hoja de Trabajo de 9 Columnas (Llamando a tu nueva función en help)
	hoja9Cols := help.HojaDeTrabajoIndustrial(balanceRaw)

	costos := help.CalcularCostosIndustriales(balanceRaw)

	DBResultados, TotResu := help.Resultados(balanceRaw, costos.CostoProduccion)

	DBBalnce, TotBalance := help.GenerarTodoElBalance(balan, TotResu.UtilidadNeta)

	// ... después de obtener TotResu y TotBalance
	d := models.DatosDashboard{
		PasivoCorriente:   TotBalance.PasivoCorriente,
		InventarioInicial: TotBalance.InventarioInicial,
		ActivoCorriente:   TotBalance.ActivoCorriente,
		InventarioFinal:   TotBalance.Inventario,
		CostoVentas:       TotResu.CostoVentas,
		ActivoTotal:       TotBalance.ActivoTotal,
		PasivoTotal:       TotBalance.PasivoTotal,
		Ventas:            TotResu.Ventas,
		UtilidadNeta:      TotResu.UtilidadNeta,
		GastosFijos:       TotResu.GastosFijos,
		GastosVariables:   TotResu.GastosVariables,
		GastosNoEfectivo:  TotResu.GastosNoEfectivo,
	}

	dbIndustrial, indicesTotales := help.GenerarDashboard(d)

	// 1. Marshalling a JSON para JS (Igual que en GenEstados)
	keysJSON, _ := json.Marshal(hoja9Cols)
	resultadosJSON, _ := json.Marshal(DBResultados)
	balanceJSON, _ := json.Marshal(DBBalnce)
	indicesJSON, _ := json.Marshal(dbIndustrial)

	// 2. Cálculo del porcentaje de margen real
	porcentajeMargen := 0.0
	if TotResu.Ventas > 0 {
		porcentajeMargen = (indicesTotales.MargenContribucion / TotResu.Ventas) * 100
	}
	resCalculados := map[string]string{
		"InvInicialMP":    config.FCont(costos.InvInicialMP),
		"ComprasMP":       config.FCont(costos.ComprasMP),
		"GastosMP":        config.FCont(costos.GastosMP),
		"DevMP":           config.FCont(costos.DevMP),
		"ComprasNetasMP":  config.FCont(costos.ComprasNetasMP),
		"MPDisponible":    config.FCont(costos.MPDisponible),
		"InvFinalMP":      config.FCont(costos.InvFinalMP),
		"MPConsumida":     config.FCont(costos.MPConsumida),
		"MOD":             config.FCont(costos.MOD),
		"CostoPrimo":      config.FCont(costos.CostoPrimo),
		"CIF":             config.FCont(costos.CIF),
		"CostoProduccion": config.FCont(costos.CostoProduccion),
		"InvInicialProc":  config.FCont(costos.InvInicialProc),
		"InvFinalProc":    config.FCont(costos.InvFinalProc),
		"CostoArtTerm":    config.FCont(costos.CostoArtTerm),
		"InvInicialPT":    config.FCont(costos.InvInicialPT),
		"InvFinalPT":      config.FCont(costos.InvFinalPT),
		"CostoVentas":     config.FCont(costos.CostoVentas),
	}
	// 3. Renderizado con nombres unificados para reutilizar HBS
	return c.Render("costos", fiber.Map{
		"Title":              "Auditoría Industrial - Costos",
		"keys":               hoja9Cols,     // Unificado
		"resultados":         DBResultados,  // Unificado
		"BalanceRows":        DBBalnce,      // Unificado
		"indices":            dbIndustrial,  // Unificado
		"res":                resCalculados, // Específico industrial
		"MontoMP":            costos.MPConsumida,
		"MontoMOD":           costos.MOD,
		"MontoCIF":           costos.CIF,
		"keysjson":           string(keysJSON),
		"resultadosjson":     string(resultadosJSON),
		"BalanceRowsjson":    string(balanceJSON),
		"indicesjson":        string(indicesJSON),
		"VentasNetas":        TotResu.VentasNetas,
		"MargenContribucion": fmt.Sprintf("%.2f", porcentajeMargen),
		"PuntoE":             indicesTotales.PuntoEContable,
		"PuntoEfe":           indicesTotales.PuntoECaja,
	})
}
