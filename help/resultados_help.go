// Package help genera las tablas de los estados
package help

import (
	"strconv"

	"ef/config"
	"ef/models"
)

func FilaCostoProduccion(costoTotal float64) ([]models.KR, float64) {
	var res []models.KR
	var dato models.ReString

	dato.Nombre = "(+) Costo de Producción"
	dato.Col2 = config.FCont(costoTotal)
	// Key especial para identificar que viene de la fábrica
	res = append(res, models.KR{Key: "399999", Value: dato})

	return res, costoTotal
}

// Ingresos financieros

func RecCol3(Balance map[string]float64, Recorido []models.Cue) ([]models.KR, float64) {
	var (
		res   []models.KR
		data  models.ReString
		Saldo float64
	)

	data.ClasNombre = "titulo"
	data.Nombre = "Ingresos"
	res = append(res, models.KR{Key: "210000", Value: data})
	for _, v := range Recorido {
		i := strconv.Itoa(v.Codigo)
		Val := Balance[i]
		if Val == 0 {
			continue
		}
		var dato models.ReString
		dato.Nombre = v.Nombre
		dato.Col3 = config.FCont(Val)
		if v.Saldo == "perdida" {
			dato.Cla3 = "total"
			Saldo -= Val
		}
		if v.Saldo == "ganancia" {
			Saldo += Val
		}
		res = append(res, models.KR{Key: i, Value: dato})
	}
	data.Col3 = config.FCont(Saldo)
	data.Cla3 = "fw-bold border-top"
	data.ClasNombre = ""
	data.Nombre = "Ventas Brutas"
	res = append(res, models.KR{Key: "210000", Value: data})

	return res, Saldo
}

func RecCol1(Balance map[string]float64, Recorido []models.Cue) ([]models.KR, float64) {
	var res []models.KR
	var Saldo float64
	for _, v := range Recorido {
		i := strconv.Itoa(v.Codigo)
		Val := Balance[i]
		if Val == 0 {
			continue
		}
		var dato models.ReString
		dato.Nombre = v.Nombre
		dato.Col1 = config.FCont(Val)
		saldo := v.Saldo
		if saldo == "perdida" {
			Saldo += Val
		} else {
			Saldo -= Val
		}
		res = append(res, models.KR{Key: i, Value: dato})
	}
	// EL TRUCO: Al último elemento le ponemos el total en la Columna 2
	if len(res) > 0 {
		res[len(res)-1].Value.Col2 = config.FCont(Saldo)
		res[len(res)-1].Value.Cla1 = "fw-bold border-bottom" // Línea de suma
	}
	return res, Saldo
}

func RecCol1Tot(Balance map[string]float64, Recorido []models.Cue, Sum float64) ([]models.KR, float64) {
	var res []models.KR
	var Saldo float64
	for _, v := range Recorido {
		i := strconv.Itoa(v.Codigo)
		Val := Balance[i]
		if Val == 0 {
			continue
		}

		var dato models.ReString
		dato.Nombre = v.Nombre
		dato.Col1 = config.FCont(Val)
		saldo := v.Saldo
		if saldo == "perdida" {
			Saldo += Val
		} else {
			Saldo -= Val
		}
		res = append(res, models.KR{Key: i, Value: dato})
	}

	// EL TRUCO: Al último elemento le ponemos el total en la Columna 2
	if len(res) > 0 {
		res[len(res)-1].Value.Col2 = config.FCont(Saldo)
		res[len(res)-1].Value.Col3 = config.FCont(Sum + Saldo)
		res[len(res)-1].Value.Cla2 = "fw-bold border-bottom"
		res[len(res)-1].Value.Cla3 = "fw-bold border-bottom"
		res[len(res)-1].Value.Cla1 = "border-bottom" // Línea de suma
	}
	return res, Saldo
}

func InvIni(Balance map[string]float64) ([]models.KR, float64) {
	var (
		res  []models.KR
		sum  float64
		dato models.ReString
	)
	dato.Nombre = "Costo de Ventas"
	dato.ClasNombre = "titulo"
	res = append(res, models.KR{Key: "220000", Value: dato})
	for _, V := range config.InveIni {
		k := strconv.Itoa(V.Codigo)
		Val := Balance[k]
		var dato models.ReString
		dato.Nombre = V.Nombre
		dato.Col2 = config.FCont(Val)
		res = append(res, models.KR{Key: k, Value: dato})
		sum += Val
	}
	return res, sum
}

func Comp(CatalogoResultados map[string]float64) ([]models.KR, float64) {
	var res []models.KR
	var sum float64

	for i, v := range config.Compras {
		k := strconv.Itoa(v.Codigo)
		Val := CatalogoResultados[k]

		// Siempre sumamos al acumulador para el cálculo técnico
		if v.Saldo == "perdida" {
			sum += Val
		} else {
			sum -= Val
		}

		// Solo dibujamos la fila si tiene valor
		if Val != 0 {
			res = append(res, models.KR{Key: k, Value: models.ReString{
				Nombre: v.Nombre,
				Col1:   config.FCont(Val),
			}})
		}

		// Insertar subtotales estructurales
		if i == (len(config.Compras) - 2) { // Compras Brutas
			res = append(res, models.KR{Key: "220200", Value: models.ReString{
				Nombre: "Compras Brutas", Col1: config.FCont(sum), Cla1: "fw-bold border-top",
			}})
		}
		if i == (len(config.Compras) - 1) { // Compras Netas
			res = append(res, models.KR{Key: "220299", Value: models.ReString{
				Nombre: "Compras Netas", Col2: config.FCont(sum), Cla1: "fw-bold border-top", Cla2: "fw-bold border-bottom",
			}})
		}
	}
	return res, sum
}

func CostoVentas(Balance map[string]float64, invIni float64, compNetas float64) ([]models.KR, float64) {
	var res []models.KR
	invFin := Balance["220005"]

	disponible := invIni + compNetas

	// Fila 1: Disponibles (Columna 2)
	res = append(res, models.KR{Key: "DISP", Value: models.ReString{
		Nombre: "Mercaderías Disponibles",
		Col2:   config.FCont(disponible),
	}})

	// Fila 2: Inventario Final (Columna 2 con resta)
	res = append(res, models.KR{Key: "220005", Value: models.ReString{
		Nombre: "(-) Inventario Final de Mercaderías",
		Col2:   config.FCont(invFin),
		Cla2:   "border-bottom",
	}})

	costoVentas := disponible - invFin

	// Fila 3: EL TOTAL (Debe ir en Columna 3 para restar a las Ventas)
	res = append(res, models.KR{Key: "CV_TOTAL", Value: models.ReString{
		Nombre: "COSTO DE VENTAS",
		Col3:   config.FCont(costoVentas),
		// Cla3:   "fw-bold border-top",
	}})

	return res, costoVentas
}

func Resultados(Balance map[string]float64, costoProd float64) ([]models.KR, models.TotalesResultados) {
	Res := make([]models.KR, 0)
	var t models.TotalesResultados

	// --- 1. INGRESOS ---
	filasIng, vVentasNetas := RecCol3(Balance, config.Ingresos) // Usamos tu variable IngresosCo
	Res = append(Res, filasIng...)

	// --- 2. COSTO DE VENTAS (Híbrido) ---
	if costoProd > 0 {
		// --- ESCENARIO INDUSTRIAL ---
		// 2.1 Inv. Inicial de Productos Terminados
		vInvIniPT := Balance["220002"] // Código que definiste en InveIniCo
		var datoIni models.ReString
		datoIni.Nombre = "Inventario Inicial de Productos Terminados"
		datoIni.Col2 = config.FCont(vInvIniPT)
		Res = append(Res, models.KR{Key: "220002", Value: datoIni})

		// 2.2 (+) Costo de Producción (Lo que viene de la fábrica)
		var datoCP models.ReString
		datoCP.Nombre = "(+) Costo de Producción"
		datoCP.Col2 = config.FCont(costoProd)
		Res = append(Res, models.KR{Key: "CP", Value: datoCP})

		// 2.3 (-) Inv. Final de Productos Terminados
		vInvFinPT := Balance["111305"] // Código que definiste en RealisableCo
		var datoFin models.ReString
		datoFin.Nombre = "(-) Inventario Final de Productos Terminados"
		datoFin.Col2 = config.FCont(vInvFinPT)
		Res = append(Res, models.KR{Key: "111305", Value: datoFin})

		t.CostoVentas = (vInvIniPT + costoProd) - vInvFinPT

		// AÑADE ESTO PARA QUE SE VEA EL TOTAL DEL COSTO EN EL REPORTE:
		Res = append(Res, models.KR{Key: "CV_IND", Value: models.ReString{
			Nombre:     "CALCULO COSTO VENTAS",
			Col3:       config.FCont(t.CostoVentas),
			ClasNombre: "titulo border-bottom",
		}})
	} else {
		// --- ESCENARIO COMERCIAL (Tu lógica original) ---
		filasInvIni, vInvIni := InvIni(Balance)
		Res = append(Res, filasInvIni...)

		filasComp, vCompNetas := Comp(Balance)
		Res = append(Res, filasComp...)

		filasCV, vCV := CostoVentas(Balance, vInvIni, vCompNetas)
		Res = append(Res, filasCV...)
		t.CostoVentas = vCV
	}

	// --- 3. MARGEN BRUTO ---
	vUtilidadBruta := vVentasNetas - t.CostoVentas
	Res = append(Res, models.KR{Key: "UB", Value: models.ReString{
		Nombre: "MARGEN BRUTO (UTILIDAD BRUTA)",
		Col3:   config.FCont(vUtilidadBruta),
		Cla3:   "fw-bold border-top",
	}})

	// --- 2. GASTOS DE OPERACIÓN ---
	Res = append(Res, models.KR{Key: "TIT_GO", Value: models.ReString{Nombre: "GASTOS DE OPERACIÓN", ClasNombre: "fw-bold titulo"}})
	Res = append(Res, models.KR{Key: "TIT_GD", Value: models.ReString{Nombre: "Gastos de Distribución", ClasNombre: "fw-bold"}})
	gtoven, vSumGtoVen := RecCol1(Balance, config.GtoVentas)
	Res = append(Res, gtoven...)
	Res = append(Res, models.KR{Key: "TIT_GA", Value: models.ReString{Nombre: "Gastos de Administración", ClasNombre: "fw-bold"}})
	gtoadm, vSumGtoAdm := RecCol1Tot(Balance, config.GtoAdmin, vSumGtoVen)
	Res = append(Res, gtoadm...)

	vTotalGastos := vSumGtoVen + vSumGtoAdm
	vResultadoOp := vUtilidadBruta - vTotalGastos

	Res = append(Res, models.KR{Key: "RO", Value: models.ReString{
		Nombre: "RESULTADO DE OPERACIÓN",
		Col3:   config.FCont(vResultadoOp),
		// Cla3:   "fw-bold border-top",
		Cla2: "fw-bold border-top",
	}})

	Res = append(Res, models.KR{Key: "TIT_OI", Value: models.ReString{Nombre: "INGRESOS FINANCIEROS", ClasNombre: "fw-bold"}})
	filasIngFin, vIngFin := RecCol1(Balance, config.IngrFina)
	Res = append(Res, filasIngFin...)

	Res = append(Res, models.KR{Key: "TIT_OI", Value: models.ReString{Nombre: "GASTOS FINANCIEROS", ClasNombre: "fw-bold"}})

	filasGtoFin, vGtoFin := RecCol1Tot(Balance, config.GastosFina, vIngFin)
	Res = append(Res, filasGtoFin...)

	vSumaIngGtoFin := vResultadoOp + -(vGtoFin + vIngFin)
	// Agregamos una fila de subtotal visual para Otros Ingresos
	Res = append(Res, models.KR{Key: "SUM_OI", Value: models.ReString{
		Nombre: "Resultado Despues de Ing y Gto Financieros",
		Col3:   config.FCont(vSumaIngGtoFin),
		Cla2:   "fw-bold border-top",
	}})

	// --- SECCIÓN DE GASTOS FINANCIEROS ---
	Res = append(Res, models.KR{Key: "TIT_GE", Value: models.ReString{Nombre: "OTROS GASTOS", ClasNombre: "fw-bold"}})

	filasOtrosGto, vOtrosGto := RecCol1(Balance, config.OtrosGtoFina)
	Res = append(Res, filasOtrosGto...)
	Res = append(Res, models.KR{Key: "TIT_GE", Value: models.ReString{Nombre: "OTROS INGRESOS", ClasNombre: "fw-bold"}})
	filasOtrosIng, vOtrosIng := RecCol1Tot(Balance, config.Otring, vOtrosGto)
	Res = append(Res, filasOtrosIng...)

	vSumaOtrosGtoIng := vSumaIngGtoFin - (vOtrosGto + vOtrosIng)

	// --- 4. CÁLCULO FINAL ---
	vUtilidadNeta := vSumaOtrosGtoIng

	Res = append(Res, models.KR{Key: "UN", Value: models.ReString{
		Nombre:      "UTILIDAD NETA DEL EJERCICIO",
		ClasNombre:  "fw-bold fs-5",
		Col3:        config.FCont(vUtilidadNeta),
		Cla3:        "fw-bold fs-5 border-double", // Necesitas definir border-double en CSS
		EsResultado: true,
	}})

	// Llenar totales para el Dashboard
	t.UtilidadNeta = vUtilidadNeta
	t.MargenBruto = vUtilidadBruta
	t.VentasNetas = vVentasNetas // <--- ¡ASEGÚRATE DE AGREGAR ESTA LÍNEA!

	return Res, t
}
