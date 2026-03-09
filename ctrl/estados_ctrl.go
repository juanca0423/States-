// Package ctrl elcontrolador para generar los estados financieros
package ctrl

import (
	//	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"ef/config"
	"ef/help"
	"ef/models"
	_ "image/jpeg" // Esto registra el de JPEG por si acaso
	_ "image/png"  // Esto registra el decodificador de PNG

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
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

	// 3. Ahora el dashboard tiene TODO lo que necesita
	dashboard := help.GenerarDashboard(
		TotBalance.ActivoCorriente,
		TotBalance.PasivoCorriente,
		TotBalance.InventarioInicial, // <--- Ahora lo tienes en el struct
		TotBalance.Inventario,        // <--- Inventario Final
		TotResu.CostoVentas,
		TotBalance.ActivoTotal,
		TotBalance.PasivoTotal,
		TotResu.VentasNetas,
		TotResu.UtilidadNeta,
	)

	type ItemBalance struct {
		Codigo string
		Monto  float64
	}
	var listaBalance []ItemBalance
	for k, v := range Balance {
		listaBalance = append(listaBalance, ItemBalance{Codigo: k, Monto: v})
	}
	// 5. Renderizado
	return c.Render("estados", fiber.Map{
		"Title":           "Estados Financieros",
		"keys":            DBHojadetravajo,
		"resultados":      DBResultados,
		"BalanceRows":     DBBalance,
		"indices":         dashboard,
		"BalanceOriginal": listaBalance, // <--- Enviamos la LISTA, no el mapa
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

	costos := help.CalcularCostosIndustriales(balanceRaw)
	// Se lo pasas a la función para que arme el Estado de Resultados Industrial

	DBResultados, TotResu := help.Resultados(balanceRaw, costos.CostoProduccion)

	DBBalnce, _ := help.GenerarTodoElBalance(balan, TotResu.UtilidadNeta)
	// ... (Resultados y Balance igual)

	// 5. Empaquetar (Agregamos los campos para las 3 columnas del HBS)
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
	// ... debajo de donde calculas 'costos'
	ventas := balanceRaw["210001"] // Código de tu nomenclatura
	dbIndustrial := help.GenerarDashboardIndustrial(costos, ventas)

	return c.Render("costos", fiber.Map{
		"Title":               "Auditoría Industrial - Costos",
		"filasHoja":           hoja9Cols,
		"Resultados":          DBResultados,
		"balanceRaw":          DBBalnce,
		"res":                 resCalculados,
		"dashboardIndustrial": dbIndustrial, // <--- Nueva variable para el HBS
	})
}

// Función auxiliar para limpiar el formato de moneda antes de meterlo a Excel
func clean(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ",", ""), "Q ", "")
}

func AgregarPestañaCierre(f *excelize.File, resultados []models.KR) {
	sheet := "Partida de Cierre"
	f.NewSheet(sheet)

	// 1. DEFINICIÓN DE ESTILOS
	estiloNegrita, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	estiloMoneda, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})

	// Estilo para el cierre: Negrita + Moneda + Línea simple arriba + Doble línea abajo
	estiloCierre, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},    // Línea simple
			{Type: "bottom", Color: "000000", Style: 6}, // Doble línea
		},
	})

	// 2. ENCABEZADOS (Sin columna de código)
	f.SetCellValue(sheet, "A1", "CUENTA / DESCRIPCIÓN")
	f.SetCellValue(sheet, "B1", "DEBE")
	f.SetCellValue(sheet, "C1", "HABER")
	f.SetCellStyle(sheet, "A1", "C1", estiloNegrita)

	f.SetCellValue(sheet, "A2", "Pda. 1 - Liquidación de Cuentas de Resultados")
	f.SetCellStyle(sheet, "A2", "A2", estiloNegrita)

	rowIdx := 3
	var sumaDebe, sumaHaber float64

	// 3. PROCESAMIENTO DE CUENTAS
	for _, r := range resultados {
		if r.Key == "UB" || r.Key == "RO" || r.Key == "UN" || r.Key == "TIT_GO" || r.Value.Nombre == "" {
			continue
		}

		valCol3, _ := strconv.ParseFloat(clean(r.Value.Col3), 64) // Ingresos
		valCol1, _ := strconv.ParseFloat(clean(r.Value.Col1), 64) // Gastos
		valCol2, _ := strconv.ParseFloat(clean(r.Value.Col2), 64) // Gastos/Subtotales

		if valCol3 != 0 {
			// INGRESO: Se carga al DEBE para cerrar
			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), r.Value.Nombre)
			f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), valCol3)
			sumaDebe += valCol3
			rowIdx++
		} else if valCol1 != 0 || valCol2 != 0 {
			// GASTO: Se abona al HABER para cerrar
			gasto := valCol1
			if gasto == 0 {
				gasto = valCol2
			}

			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "    "+r.Value.Nombre) // Sangría
			f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), gasto)
			sumaHaber += gasto
			rowIdx++
		}
	}

	// 4. UTILIDAD PARA CUADRAR
	utilidad := sumaDebe - sumaHaber
	f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "    UTILIDAD NETA DEL EJERCICIO")
	f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), utilidad)
	sumaHaber += utilidad
	rowIdx++

	// 5. FILA FINAL DE SUMAS IGUALES
	f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "SUMAS IGUALES")
	f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), sumaDebe)
	f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), sumaHaber)

	// 6. APLICACIÓN FINAL DE ESTILOS Y FORMATOS
	// Formato moneda a todos los datos de las columnas B y C
	f.SetCellStyle(sheet, "B3", "C"+strconv.Itoa(rowIdx-1), estiloMoneda)

	// Estilo de cierre (Doble línea) a los totales
	f.SetCellStyle(sheet, "B"+strconv.Itoa(rowIdx), "C"+strconv.Itoa(rowIdx), estiloCierre)
	f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloNegrita)

	// Ajuste de anchos
	f.SetColWidth(sheet, "A", "A", 50)
	f.SetColWidth(sheet, "B", "C", 18)
}

func AgregarPestañaCierreBalance(f *excelize.File, balanceRows []models.BaString, utilidad float64) {
	sheet := "Cierre de Balance"
	f.NewSheet(sheet)

	estiloNegrita, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	estiloMoneda, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})

	// Estilo combinado: Moneda + Negrita + Bordes Contables
	estiloTotal, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		// Formato de moneda con comas y dos decimales
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},    // Línea simple arriba
			{Type: "bottom", Color: "000000", Style: 6}, // Línea doble abajo (Style 6 es doble)
		},
	})

	f.SetCellValue(sheet, "A1", "DESCRIPCIÓN DE CUENTA")
	f.SetCellValue(sheet, "B1", "DEBE")
	f.SetCellValue(sheet, "C1", "HABER")
	f.SetCellStyle(sheet, "A1", "C1", estiloNegrita)

	rowIdx := 3
	var sumaDebe, sumaHaber float64

	for _, row := range balanceRows {
		valor, _ := strconv.ParseFloat(clean(row.Col1), 64)
		if valor == 0 {
			continue
		}

		cod := strings.TrimSpace(row.Codigo)
		nombre := strings.ToLower(row.Nombre)
		esActivo := strings.HasPrefix(cod, "11")
		esPasivoOCapital := strings.HasPrefix(cod, "12") || strings.HasPrefix(cod, "13")
		restaEnSuSeccion := strings.Contains(nombre, "(-)") || strings.Contains(nombre, "dep. acu.") ||
			strings.Contains(nombre, "amor. acu.") || strings.Contains(nombre, "res. para")

		if (esActivo && !restaEnSuSeccion) || (esPasivoOCapital && restaEnSuSeccion) {
			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "    "+row.Nombre)
			f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), math.Abs(valor))
			sumaHaber += math.Abs(valor)
		} else {
			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), row.Nombre)
			f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), math.Abs(valor))
			sumaDebe += math.Abs(valor)
		}
		rowIdx++
	}

	if utilidad != 0 {
		f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "Utilidad del Ejercicio")
		f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), utilidad)
		sumaDebe += utilidad
		rowIdx++
	}

	f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), "SUMAS IGUALES")
	f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), sumaDebe)
	f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), sumaHaber)

	// APLICACIÓN DE ESTILOS SIN SOBREPOSICIÓN
	f.SetCellStyle(sheet, "B3", "C"+strconv.Itoa(rowIdx-1), estiloMoneda)                    // Datos normales
	f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloNegrita) // Etiqueta "SUMAS IGUALES"
	f.SetCellStyle(sheet, "B"+strconv.Itoa(rowIdx), "C"+strconv.Itoa(rowIdx), estiloTotal)   // Totales (Moneda + Negrita)

	f.SetColWidth(sheet, "A", "A", 45)
	f.SetColWidth(sheet, "B", "C", 18)
}

func AgregarCaratula(f *excelize.File, nombreEmpresa string) {
	sheet := "Carátula"
	f.NewSheet(sheet)
	f.SetSheetName("Sheet1", sheet)

	estiloEmpresa, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 24, Color: "000080"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// 1. Insertar un SOLO Logo (Escalado al 0.1 como pediste)
	mantenerAspecto := true
	_ = f.AddPicture(sheet, "B2", "./estaticos/images/logo.png", &excelize.GraphicOptions{
		ScaleX:          0.1,
		ScaleY:          0.1,
		LockAspectRatio: mantenerAspecto,
	})

	// 2. Nombre de la Empresa y Títulos bajados a la fila 10
	f.MergeCell(sheet, "B10", "H10")
	f.SetCellValue(sheet, "B10", strings.ToUpper(nombreEmpresa))
	f.SetCellStyle(sheet, "B10", "B10", estiloEmpresa)

	f.MergeCell(sheet, "B11", "H11")
	f.SetCellValue(sheet, "B11", "ESTADOS FINANCIEROS Y CIERRE CONTABLE")

	f.MergeCell(sheet, "B12", "H12")
	f.SetCellValue(sheet, "B12", "Período: Ejercicio Fiscal 2026")

	// 3. Info del sistema
	f.SetCellValue(sheet, "C15", "Generado por:")
	f.SetCellValue(sheet, "D15", "Sistema Contable Gemini AI")

	// Quitar líneas de cuadrícula
	showGrid := false
	_ = f.SetSheetView(sheet, 0, &excelize.ViewOptions{
		ShowGridLines: &showGrid,
	})
}

func AgregarPestañaResumen(f *excelize.File, indicadores []models.IndiceFinanciero) {
	sheet := "Resumen Ejecutivo"
	f.NewSheet(sheet)

	// --- ESTILOS ---
	estiloTitulo, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14, Color: "#4A86E8"},
	})
	estiloHeaderTabla, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4A86E8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	// Estilo para números con 2 decimales
	estiloNumero, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &[]string{"0.00"}[0],
		Alignment:    &excelize.Alignment{Horizontal: "center"},
	})

	f.SetCellValue(sheet, "A1", "RESUMEN DE INDICADORES FINANCIEROS")
	f.SetCellStyle(sheet, "A1", "A1", estiloTitulo)

	headers := []string{"INDICADOR", "VALOR", "INTERPRETACIÓN", "FÓRMULA / DETALLE"}
	for i, h := range headers {
		celda, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheet, celda, h)
		f.SetCellStyle(sheet, celda, celda, estiloHeaderTabla)
	}

	rowIdx := 4
	for _, ind := range indicadores {
		f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), ind.Nombre)

		// --- CONVERSIÓN DE TEXTO A NÚMERO ---
		// Limpiamos el string de símbolos como Q, $, % o comas
		valorLimpio := strings.NewReplacer("Q", "", "$", "", "%", "", ",", "", " ", "").Replace(ind.Valor)

		if v, err := strconv.ParseFloat(valorLimpio, 64); err == nil {
			// Si la conversión es exitosa, lo mandamos como número
			f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), v)
			f.SetCellStyle(sheet, "B"+strconv.Itoa(rowIdx), "B"+strconv.Itoa(rowIdx), estiloNumero)
		} else {
			// Si falla (ej: un texto), lo mandamos como string original
			f.SetCellValue(sheet, "B"+strconv.Itoa(rowIdx), ind.Valor)
		}

		f.SetCellValue(sheet, "C"+strconv.Itoa(rowIdx), ind.Interpretacion)
		f.SetCellValue(sheet, "D"+strconv.Itoa(rowIdx), ind.DetalleCuenta)

		// --- COLORES SEGÚN LA CLASE ---
		if ind.Clase != "" {
			var color string
			switch ind.Clase {
			case "text-success":
				color = "#C6EFCE"
			case "text-danger":
				color = "#FFC7CE"
			case "text-warning":
				color = "#FFEB9C"
			}
			if color != "" {
				estiloColor, _ := f.NewStyle(&excelize.Style{
					Fill:         excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
					Font:         &excelize.Font{Bold: true},
					Alignment:    &excelize.Alignment{Horizontal: "center"},
					CustomNumFmt: &[]string{"0.00"}[0],
				})
				f.SetCellStyle(sheet, "B"+strconv.Itoa(rowIdx), "B"+strconv.Itoa(rowIdx), estiloColor)
			}
		}
		rowIdx++
	}

	f.SetColWidth(sheet, "A", "A", 35)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "D", 55)
}

func AgregarPestañaNotas(f *excelize.File) {
	sheet := "Notas a los EEFF"
	f.NewSheet(sheet)

	// Estilo para títulos de notas
	estiloNotaTit, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Color: "#FFFFFF"}, Fill: excelize.Fill{Type: "pattern", Color: []string{"#4A86E8"}, Pattern: 1}})
	estiloTexto, _ := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"}})

	notas := []struct {
		Titulo      string
		Placeholder string
	}{
		{"NOTA 1. IDENTIFICACIÓN DE LA ENTIDAD", "La empresa [Nombre de la Empresa] fue constituida según escritura pública... Su actividad principal es..."},
		{"NOTA 2. BASES DE PRESENTACIÓN", "Los estados financieros se prepararon de acuerdo con Normas Internacionales de Información Financiera (NIIF) para PyMEs..."},
		{"NOTA 3. POLÍTICAS CONTABLES SIGNIFICATIVAS", "Efectivo: Se registra al valor nominal. Inventarios: Se valúan bajo el método [PEPS/Promedio]."},
		{"NOTA 4. EFECTIVO Y EQUIVALENTES", "El saldo representa disponibilidades en cuentas monetarias y ahorros en moneda local..."},
		{"NOTA 5. PROPIEDAD, PLANTA Y EQUIPO", "Los activos se registran al costo de adquisición y se deprecian bajo el método de línea recta..."},
	}

	rowIdx := 1
	for _, n := range notas {
		// Escribir Título
		f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), n.Titulo)
		f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloNotaTit)

		// Escribir el Placeholder
		rowIdx++
		f.MergeCell(sheet, "A"+strconv.Itoa(rowIdx), "H"+strconv.Itoa(rowIdx+2)) // Espacio grande para escribir
		f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), n.Placeholder)
		f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloTexto)

		rowIdx += 4 // Espacio entre notas
	}

	f.SetColWidth(sheet, "A", "A", 100)
}

func AgregarPestañaDashboard(f *excelize.File, indices []models.IndiceFinanciero) {
	sheet := "Analisis"
	f.NewSheet(sheet)

	// 1. ESTILOS PRE-DEFINIDOS
	// Estilo Normal (Blanco/Estándar)
	estiloNormal, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Estilo de Alerta (Rojo) - Aplicado directamente por Go
	estiloAlerta, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFC7CE"}, Pattern: 1},
		Font:      &excelize.Font{Color: "#9C0006"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	f.SetCellValue(sheet, "A1", "INDICADOR")
	f.SetCellValue(sheet, "B1", "VALOR")
	f.SetColWidth(sheet, "A", "A", 35)
	f.SetColWidth(sheet, "B", "B", 15)

	if len(indices) == 0 {
		return
	}

	var builder strings.Builder
	for i, ind := range indices {
		row := i + 2
		builder.Reset()

		for _, char := range ind.Valor {
			if (char >= '0' && char <= '9') || char == '.' || char == '-' {
				builder.WriteRune(char)
			}
		}

		if builder.Len() > 0 {
			num, _ := strconv.ParseFloat(builder.String(), 64)
			celdaValor := "B" + strconv.Itoa(row)
			celdaNombre := "A" + strconv.Itoa(row)

			f.SetCellValue(sheet, celdaNombre, ind.Nombre)
			f.SetCellValue(sheet, celdaValor, num)

			// --- EL "FORMATO CONDICIONAL" HECHO POR GO ---
			// Evaluamos aquí la lógica. Si es menor a 1, aplicamos el estilo rojo.
			if num < 1 {
				f.SetCellStyle(sheet, celdaValor, celdaValor, estiloAlerta)
			} else {
				f.SetCellStyle(sheet, celdaValor, celdaValor, estiloNormal)
			}
		}
	}

	// 2. ELIMINAMOS POR COMPLETO f.SetConditionalFormat
	// Esto garantiza que el XML de la sheet4 sea 100% estándar y limpio.

	// 3. GRÁFICO (Opcional, si quieres mantenerlo)
	lastRow := len(indices) + 1
	grafico := &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Valor",
				Categories: fmt.Sprintf("'%s'!$A$2:$A$%d", sheet, lastRow),
				Values:     fmt.Sprintf("'%s'!$B$2:$B$%d", sheet, lastRow),
			},
		},
		Title: []excelize.RichTextRun{{Text: "KPIs Financieros"}},
	}
	_ = f.AddChart(sheet, "E1", grafico)
}

func AgregarPestañaHojaTrabajo(f *excelize.File, datos []models.KV) {
	sheet := "Hoja de Trabajo"
	f.NewSheet(sheet)

	// --- COLORES Y ESTILOS (Tu paleta) ---
	azulSuave := "#4A86E8"
	grisClaro := "#F3F3F3"
	verdeSuave := "#E2F0D9"
	celesteClaro := "#DDEBF7"
	beigeClaro := "#FFF2CC"

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{azulSuave}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Estilos de bloque con formato numérico incorporado
	estiloSaldos, _ := f.NewStyle(&excelize.Style{
		Fill:         excelize.Fill{Type: "pattern", Color: []string{grisClaro}, Pattern: 1},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})
	estiloResultados, _ := f.NewStyle(&excelize.Style{
		Fill:         excelize.Fill{Type: "pattern", Color: []string{verdeSuave}, Pattern: 1},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})
	estiloBalance, _ := f.NewStyle(&excelize.Style{
		Fill:         excelize.Fill{Type: "pattern", Color: []string{celesteClaro}, Pattern: 1},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})

	// --- ENCABEZADOS ---
	f.SetCellValue(sheet, "A1", "DESCRIPCIÓN DE LA CUENTA")
	f.MergeCell(sheet, "B1", "C1")
	f.SetCellValue(sheet, "B1", "BALANCE DE SALDOS")
	f.MergeCell(sheet, "D1", "E1")
	f.SetCellValue(sheet, "D1", "ESTADO DE RESULTADOS")
	f.MergeCell(sheet, "F1", "G1")
	f.SetCellValue(sheet, "F1", "BALANCE GENERAL")

	f.SetCellValue(sheet, "B2", "DEBE")
	f.SetCellValue(sheet, "C2", "HABER")
	f.SetCellValue(sheet, "D2", "PÉRDIDA")
	f.SetCellValue(sheet, "E2", "GANANCIA")
	f.SetCellValue(sheet, "F2", "ACTIVO")
	f.SetCellValue(sheet, "G2", "PASIVO")
	f.SetCellStyle(sheet, "A1", "G2", headerStyle)

	// --- LLENADO DE DATOS ---
	rowIdx := 3
	for i := range datos {
		// Acceso directo por índice para evitar copias de valor vacías
		ht := datos[i].Value

		// Forzar escritura del nombre para verificar
		f.SetCellValue(sheet, "A"+strconv.Itoa(rowIdx), ht.Nombre)

		// Aplicar fondo beige a filas pares para el nombre
		if i%2 == 0 {
			estiloBeige, _ := f.NewStyle(&excelize.Style{Fill: excelize.Fill{Type: "pattern", Color: []string{beigeClaro}, Pattern: 1}})
			f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloBeige)
		}

		// Mapeo de columnas
		columnas := []struct {
			val   string
			letra string
			style int
		}{
			{ht.Debe, "B", estiloSaldos},
			{ht.Haber, "C", estiloSaldos},
			{ht.Perdidas, "D", estiloResultados},
			{ht.Ganancias, "E", estiloResultados},
			{ht.Activo, "F", estiloBalance},
			{ht.Pasivo, "G", estiloBalance},
		}

		for _, col := range columnas {
			// Cambiamos clean(col.val) por una limpieza manual rápida si clean falla
			vStr := strings.ReplaceAll(strings.ReplaceAll(col.val, ",", ""), " ", "")
			if v, err := strconv.ParseFloat(vStr, 64); err == nil && v != 0 {
				celda := col.letra + strconv.Itoa(rowIdx)
				f.SetCellValue(sheet, celda, v)
				f.SetCellStyle(sheet, celda, celda, col.style)
			}
		}

		// Totales
		nombreUpper := strings.ToUpper(ht.Nombre)
		if strings.Contains(nombreUpper, "SUMAS") || strings.Contains(nombreUpper, "GANANCIA") || strings.Contains(nombreUpper, "PÉRDIDA") {
			totalStyle, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{Bold: true},
				Border: []excelize.Border{
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 6},
				},
			})
			f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIdx), "G"+strconv.Itoa(rowIdx), totalStyle)
		}
		rowIdx++
	}

	// --- AJUSTES DE VISTA ---
	f.SetColWidth(sheet, "A", "A", 42)
	f.SetColWidth(sheet, "B", "G", 15)

	// Si esto sigue dando problemas, comenta estas 3 líneas para probar:
	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      1,
		YSplit:      2,
		TopLeftCell: "B3",
	})
}

func ExportarExcelTodo(c *fiber.Ctx) error {
	balanceMap := make(map[string]float64)
	c.Context().PostArgs().VisitAll(func(key, value []byte) {
		if val, err := strconv.ParseFloat(string(value), 64); err == nil {
			if val != 0 {
				balanceMap[string(key)] = val
			}
		}
	})

	balan := config.DividirCuentas(balanceMap)
	dbResultados, totResu := help.Resultados(balanceMap, 0)
	dbBalance, totBalance := help.GenerarTodoElBalance(balan, totResu.UtilidadNeta)
	// 3. Ahora el dashboard tiene TODO lo que necesita
	dashboard := help.GenerarDashboard(
		totBalance.ActivoCorriente,
		totBalance.PasivoCorriente,
		totBalance.InventarioInicial, // <--- Ahora lo tienes en el struct
		totBalance.Inventario,        // <--- Inventario Final
		totResu.CostoVentas,
		totBalance.ActivoTotal,
		totBalance.PasivoTotal,
		totResu.VentasNetas,
		totResu.UtilidadNeta,
	)

	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "Carátula")
	AgregarCaratula(f, "NOMBRE DE TU EMPRESA")

	// --- COLORES DE LA PALETA ---
	azulSuave := "#4A86E8"
	verdeSuave := "#E2F0D9"
	celesteClaro := "#DDEBF7"
	beigeClaro := "#FFF2CC"

	// --- ESTILOS COMPARTIDOS ---
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{azulSuave}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	estiloBold, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	estiloBeige, _ := f.NewStyle(&excelize.Style{Fill: excelize.Fill{Type: "pattern", Color: []string{beigeClaro}, Pattern: 1}})

	// 2. ESTADO DE RESULTADOS
	sheetRes := "Estado de Resultados"
	f.NewSheet(sheetRes)
	f.SetCellValue(sheetRes, "A1", "DESCRIPCIÓN")
	f.SetCellValue(sheetRes, "B1", "PARCIAL")
	f.SetCellValue(sheetRes, "C1", "SUBTOTAL")
	f.SetCellValue(sheetRes, "D1", "TOTAL")
	f.SetCellStyle(sheetRes, "A1", "D1", headerStyle)

	estiloResBase, _ := f.NewStyle(&excelize.Style{
		Fill:         excelize.Fill{Type: "pattern", Color: []string{verdeSuave}, Pattern: 1},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})

	rowIdx := 2
	for i, row := range dbResultados {
		f.SetCellValue(sheetRes, "A"+strconv.Itoa(rowIdx), row.Value.Nombre)
		if i%2 == 0 {
			f.SetCellStyle(sheetRes, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloBeige)
		}

		nombreRes := strings.ToUpper(row.Value.Nombre)
		cols := []struct{ val, let string }{{row.Value.Col1, "B"}, {row.Value.Col2, "C"}, {row.Value.Col3, "D"}}

		for _, c := range cols {
			if v, err := strconv.ParseFloat(clean(c.val), 64); err == nil && v != 0 {
				f.SetCellValue(sheetRes, c.let+strconv.Itoa(rowIdx), v)
				f.SetCellStyle(sheetRes, c.let+strconv.Itoa(rowIdx), c.let+strconv.Itoa(rowIdx), estiloResBase)
			}
		}
		if strings.Contains(nombreRes, "TOTAL") || strings.Contains(nombreRes, "UTILIDAD") {
			f.SetCellStyle(sheetRes, "A"+strconv.Itoa(rowIdx), "A"+strconv.Itoa(rowIdx), estiloBold)
		}
		rowIdx++
	}
	f.SetColWidth(sheetRes, "A", "A", 45)
	f.SetColWidth(sheetRes, "B", "D", 15)

	// 3. BALANCE GENERAL
	sheetBal := "Balance General"
	f.NewSheet(sheetBal)
	f.SetCellValue(sheetBal, "A1", "Descripción")
	f.SetCellValue(sheetBal, "A1", "Descripción")
	f.SetCellValue(sheetBal, "B1", "Parciales")
	f.SetCellValue(sheetBal, "C1", "Subtotales")
	f.SetCellValue(sheetBal, "D1", "Totales")
	f.SetCellValue(sheetBal, "E1", "Suma Total")
	f.SetCellStyle(sheetBal, "A1", "E1", headerStyle)

	estiloBalBase, _ := f.NewStyle(&excelize.Style{
		Fill:         excelize.Fill{Type: "pattern", Color: []string{celesteClaro}, Pattern: 1},
		CustomNumFmt: &[]string{"_ * #,##0.00_ ;_ * -#,##0.00_ ;_ * \"-\"??_ ;_ @_ "}[0],
	})

	rowIdxBal := 2
	for i, row := range dbBalance {
		f.SetCellValue(sheetBal, "A"+strconv.Itoa(rowIdxBal), row.Nombre)
		if i%2 == 0 {
			f.SetCellStyle(sheetBal, "A"+strconv.Itoa(rowIdxBal), "A"+strconv.Itoa(rowIdxBal), estiloBeige)
		}

		vals := []string{row.Col1, row.Col2, row.Col3, row.Col4}
		lets := []string{"B", "C", "D", "E"}

		for j, vStr := range vals {
			if v, err := strconv.ParseFloat(clean(vStr), 64); err == nil && v != 0 {
				f.SetCellValue(sheetBal, lets[j]+strconv.Itoa(rowIdxBal), v)
				f.SetCellStyle(sheetBal, lets[j]+strconv.Itoa(rowIdxBal), lets[j]+strconv.Itoa(rowIdxBal), estiloBalBase)
			}
		}
		if strings.Contains(strings.ToUpper(row.Nombre), "TOTAL") || row.Col4 != "" {
			f.SetCellStyle(sheetBal, "A"+strconv.Itoa(rowIdxBal), "A"+strconv.Itoa(rowIdxBal), estiloBold)
		}
		rowIdxBal++
	}
	f.SetColWidth(sheetBal, "A", "A", 45)
	f.SetColWidth(sheetBal, "B", "E", 15)

	// 4. PESTAÑAS ADICIONALES
	AgregarPestañaDashboard(f, dashboard)
	AgregarPestañaCierre(f, dbResultados)
	AgregarPestañaCierreBalance(f, dbBalance, totResu.UtilidadNeta)

	dbHT := help.HojaDeTrabajo(balanceMap)
	AgregarPestañaHojaTrabajo(f, dbHT) // Esta ya tiene la paleta que definimos antes
	AgregarPestañaResumen(f, dashboard)
	AgregarPestañaNotas(f)
	f.SetActiveSheet(0)
	buffer, _ := f.WriteToBuffer()
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=Auditoria_Contable_2026.xlsx")
	return c.Send(buffer.Bytes())
}
