// Package ctrl elcontrolador para generar los estados financieros
package ctrl

import (
	"encoding/json"
	"strconv"
	"strings"

	"ef/config"
	"ef/help"
	"ef/models"
	_ "image/jpeg" // Esto registra el de JPEG por si acaso
	_ "image/png"  // Esto registra el decodificador de PNG

	"github.com/gofiber/fiber/v2"
)

type KV struct {
	Key   string
	Value models.HtString
}

func GenEstados(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error al obtener el formulario: " + err.Error())
	}

	Balance := make(map[string]float64)
	// 3. Iteramos sobre form.Value (donde vv es un []string)
	for k, vv := range form.Value {
		if len(vv) > 0 && vv[0] != "" {
			valStr := strings.ReplaceAll(vv[0], ",", "")
			valStr = strings.TrimSpace(valStr)

			if f, err := strconv.ParseFloat(valStr, 64); err == nil {
				// Solo agregamos si el monto es distinto de cero
				if f != 0 {
					Balance[k] = f
				}
			}
		}
	}
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
	}

	// Llamada limpia
	dashboard := help.GenerarDashboard(d)

	type ItemBalance struct {
		Codigo string
		Monto  float64
	}
	var listaBalance []ItemBalance
	for k, v := range Balance {
		listaBalance = append(listaBalance, ItemBalance{Codigo: k, Monto: v})
	}

	// Convertimos a JSON string para que JS no se confunda
	keysJSON, _ := json.Marshal(DBHojadetravajo)
	resultadosJSON, _ := json.Marshal(DBResultados)
	balanceJSON, _ := json.Marshal(DBBalance)
	indicesJSON, _ := json.Marshal(dashboard)

	// 5. Renderizado
	return c.Render("estados", fiber.Map{
		"Title":           "Estados Financieros",
		"keys":            DBHojadetravajo,
		"resultados":      DBResultados,
		"BalanceRows":     DBBalance,
		"indices":         dashboard,
		"nombreEmpresa":   "JC Dev Systems",
		"keysjson":        string(keysJSON),
		"resultadosjson":  string(resultadosJSON),
		"BalanceRowsjson": string(balanceJSON),
		"indicesjson":     string(indicesJSON),
		"VentasNetas":     TotResu.VentasNetas,
		//"PuntoE":
		"CostosFijos": TotResu.GastosFijos,
		"CostosVar":   TotResu.GastosVariables,
	})
}

// func GenCostoProduccion(c *fiber.Ctx) error {
// 	// 1. Capturar datos del formulario industrial
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Error en el formulario")
// 	}
//
// 	// 2. Convertir y Limpiar datos
// 	balanceRaw := make(map[string]float64)
// 	for k, vv := range form.Value {
// 		if len(vv) > 0 {
// 			if f, err := strconv.ParseFloat(vv[0], 64); err == nil {
// 				balanceRaw[k] = f
// 			}
// 		}
// 	}
// 	balan := config.DividirCuentas(balanceRaw)
//
// 	// 3. Generar la Hoja de Trabajo de 9 Columnas (Llamando a tu nueva función en help)
// 	hoja9Cols := help.HojaDeTrabajoIndustrial(balanceRaw)
//
// 	costos := help.CalcularCostosIndustriales(balanceRaw)
// 	// Se lo pasas a la función para que arme el Estado de Resultados Industrial
//
// 	DBResultados, TotResu := help.Resultados(balanceRaw, costos.CostoProduccion)
//
// 	DBBalnce, _ := help.GenerarTodoElBalance(balan, TotResu.UtilidadNeta)
// 	// ... (Resultados y Balance igual)
//
// 	// 5. Empaquetar (Agregamos los campos para las 3 columnas del HBS)
// 	resCalculados := map[string]string{
// 		"InvInicialMP":    config.FCont(costos.InvInicialMP),
// 		"ComprasMP":       config.FCont(costos.ComprasMP),
// 		"GastosMP":        config.FCont(costos.GastosMP),
// 		"DevMP":           config.FCont(costos.DevMP),
// 		"ComprasNetasMP":  config.FCont(costos.ComprasNetasMP),
// 		"MPDisponible":    config.FCont(costos.MPDisponible),
// 		"InvFinalMP":      config.FCont(costos.InvFinalMP),
// 		"MPConsumida":     config.FCont(costos.MPConsumida),
// 		"MOD":             config.FCont(costos.MOD),
// 		"CostoPrimo":      config.FCont(costos.CostoPrimo),
// 		"CIF":             config.FCont(costos.CIF),
// 		"CostoProduccion": config.FCont(costos.CostoProduccion),
// 		"InvInicialProc":  config.FCont(costos.InvInicialProc),
// 		"InvFinalProc":    config.FCont(costos.InvFinalProc),
// 		"CostoArtTerm":    config.FCont(costos.CostoArtTerm),
// 		"InvInicialPT":    config.FCont(costos.InvInicialPT),
// 		"InvFinalPT":      config.FCont(costos.InvFinalPT),
// 		"CostoVentas":     config.FCont(costos.CostoVentas),
// 	}
// 	// ... debajo de donde calculas 'costos'
// 	ventas := balanceRaw["210001"] // Código de tu nomenclatura
// //	dbIndustrial := help.GenerarDashboardIndustrial(costos, ventas)
//
// 	return c.Render("costos", fiber.Map{
// 		"Title":               "Auditoría Industrial - Costos",
// 		"filasHoja":           hoja9Cols,
// 		"Resultados":          DBResultados,
// 		"balanceRaw":          DBBalnce,
// 		"res":                 resCalculados,
// 		"dashboardIndustrial": dbIndustrial, // <--- Nueva variable para el HBS
// 	})
// }
