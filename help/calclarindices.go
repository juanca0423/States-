package help

import (
	"fmt"

	"ef/config"
	"ef/models"
)

func GenerarDashboard(d models.DatosDashboard) ([]models.IndiceFinanciero, models.IndicesTotales) {
	var (
		dashboard []models.IndiceFinanciero
		totales   models.IndicesTotales
	)

	// Variables internas para cálculos complejos
	patrimonio := d.ActivoTotal - d.PasivoTotal
	utilidadBruta := d.Ventas - d.CostoVentas

	// --- 1. LIQUIDEZ CORRIENTE ---
	liq := 0.0
	if d.PasivoCorriente > 0 {
		liq = d.ActivoCorriente / d.PasivoCorriente
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Liquidez Corriente",
		Valor:          fmt.Sprintf("%.2f", liq),
		Interpretacion: "Capacidad de cubrir deudas a corto plazo.",
		Clase:          evaluarIndice(liq, 1.0, 2.0),
		DetalleCuenta:  fmt.Sprintf("Act. Corr (%s) / Pas. Corr (%s)", config.FCont(d.ActivoCorriente), config.FCont(d.PasivoCorriente)),
	})

	// --- 2. PRUEBA DEL ÁCIDO ---
	acido := 0.0
	if d.PasivoCorriente > 0 {
		acido = (d.ActivoCorriente - d.InventarioFinal) / d.PasivoCorriente
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Prueba del Ácido",
		Valor:          fmt.Sprintf("%.2f", acido),
		Interpretacion: "Liquidez inmediata (sin depender de inventarios).",
		Clase:          evaluarIndice(acido, 0.8, 1.2),
		DetalleCuenta:  fmt.Sprintf("(Act. Corr (%s)- Inv. Final (%s)) / Pas. Corr (%s) = (%s)", config.FCont(d.ActivoCorriente), config.FCont(d.InventarioFinal), config.FCont(d.PasivoCorriente), config.FCont(acido)),
	})

	// --- 3. CAPITAL DE TRABAJO (Nuevo) ---
	capTrabajo := d.ActivoCorriente - d.PasivoCorriente
	claseCap := "text-success"
	if capTrabajo < 0 {
		claseCap = "text-danger"
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Capital de Trabajo",
		Valor:          config.FCont(capTrabajo),
		Interpretacion: "Dinero disponible para operar el día a día.",
		Clase:          claseCap,
		DetalleCuenta:  fmt.Sprintf("Activo Corriente (%s) - Pasivo Corriente(%s) = (%s)", config.FCont(d.ActivoCorriente), config.FCont(d.PasivoCorriente), config.FCont(capTrabajo)),
	})

	// --- 4. ENDEUDAMIENTO TOTAL ---
	end := 0.0
	if d.ActivoTotal > 0 {
		end = (d.PasivoTotal / d.ActivoTotal) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Nivel de Endeudamiento",
		Valor:          fmt.Sprintf("%.2f%%", end),
		Interpretacion: "Porcentaje de activos financiados por terceros.",
		Clase:          evaluarIndiceEndeudamiento(end),
		DetalleCuenta:  fmt.Sprintf("(Pasivo Total (%s) / Activo Total (%s)) * 100 = (%s)", config.FCont(d.PasivoTotal), config.FCont(d.ActivoTotal), config.FCont(end)),
	})

	// --- 5. MARGEN BRUTO (Nuevo) ---
	mBruto := 0.0
	if d.Ventas > 0 {
		mBruto = (utilidadBruta / d.Ventas) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Margen Bruto",
		Valor:          fmt.Sprintf("%.2f%%", mBruto),
		Interpretacion: "Ganancia directa sobre las ventas.",
		Clase:          "text-info",
		DetalleCuenta:  fmt.Sprintf("(Utilidad Bruta (%s) / Ventas (%s)) * 100 = %.2f%%", config.FCont(utilidadBruta), config.FCont(d.Ventas), mBruto),
	})

	// --- 6. MARGEN NETO ---
	mNeto := 0.0
	if d.Ventas > 0 {
		mNeto = (d.UtilidadNeta / d.Ventas) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Margen Neto de Utilidad",
		Valor:          fmt.Sprintf("%.2f%%", mNeto),
		Interpretacion: "Ganancia real por cada quetzal vendido.",
		Clase:          "text-success",
		DetalleCuenta:  fmt.Sprintf("(Utilidad Neta (%s) / Ventas Totales (%s)) * 100 = (%.2f%%)", config.FCont(d.UtilidadNeta), config.FCont(d.Ventas), mNeto),
	})

	// --- 7. ROA (Rentabilidad sobre Activos - Nuevo) ---
	roa := 0.0
	if d.ActivoTotal > 0 {
		roa = (d.UtilidadNeta / d.ActivoTotal) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Rentabilidad s/ Activos (ROA)",
		Valor:          fmt.Sprintf("%.2f%%", roa),
		Interpretacion: "Eficacia de los activos para generar utilidad.",
		Clase:          "text-success",
		DetalleCuenta:  fmt.Sprintf("(Utilidad Neta (%s) / Activo Total (%s)) * 100 = (%.2f%%)", config.FCont(d.UtilidadNeta), config.FCont(d.ActivoTotal), roa),
	})

	// --- 8. ROE (Rentabilidad sobre Patrimonio - Nuevo) ---
	roe := 0.0
	if patrimonio > 0 {
		roe = (d.UtilidadNeta / patrimonio) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Rentabilidad s/ Patrimonio (ROE)",
		Valor:          fmt.Sprintf("%.2f%%", roe),
		Interpretacion: "Rendimiento del capital invertido por socios.",
		Clase:          "text-success",
		DetalleCuenta:  fmt.Sprintf("(Utilidad Neta (%s) / Patrimonio Neto (%s)) * 100 = %.2f%%", config.FCont(d.UtilidadNeta), config.FCont(patrimonio), roe),
	})

	// --- 9. PUNTO DE EQUILIBRIO ---
	margenContribucion := d.Ventas - d.GastosVariables
	indiceMargen := 0.0
	if d.Ventas > 0 {
		indiceMargen = margenContribucion / d.Ventas
	}

	if indiceMargen > 0 {
		totales.PuntoEContable = d.GastosFijos / indiceMargen
		// PE de Caja: Fijos menos los que no mueven efectivo (depreciaciones)
		fijosEfectivos := d.GastosFijos - d.GastosNoEfectivo
		totales.PuntoECaja = fijosEfectivos / indiceMargen
	}

	totales.MargenContribucion = margenContribucion
	totales.CostosTotales = d.GastosVariables + d.GastosFijos

	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Punto de Equilibrio (Contable)",
		Valor:          config.FCont(totales.PuntoEContable),
		Interpretacion: "Ventas necesarias para cubrir todos los costos y gastos.",
		Clase:          "text-primary fw-bold",
		DetalleCuenta:  fmt.Sprintf("Fijos (%s) / Margen Cont. (%.2f%%)", config.FCont(d.GastosFijos), indiceMargen*100),
	})

	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Mínimo Vital de Caja",
		Valor:          config.FCont(totales.PuntoECaja),
		Interpretacion: "Ventas necesarias para cubrir solo gastos que requieren efectivo.",
		Clase:          "text-warning fw-bold",
		DetalleCuenta:  fmt.Sprintf("(Fijos %s - No Efectivo %s) / Margen Cont.", config.FCont(d.GastosFijos), config.FCont(d.GastosNoEfectivo)),
	})

	// --- 10. ROTACIÓN DE INVENTARIOS (Usando promedio) ---
	promedioInv := (d.InventarioInicial + d.InventarioFinal) / 2
	rot := 0.0
	if promedioInv > 0 {
		rot = d.CostoVentas / promedioInv
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Rotación de Inventario",
		Valor:          fmt.Sprintf("%.1f veces", rot),
		Interpretacion: "Veces que el inventario se renovó en el año.",
		Clase:          "text-info",
		DetalleCuenta:  fmt.Sprintf("Costo Ventas (%s) / Inv. Promedio (%s) = %.1f veces", config.FCont(d.CostoVentas), config.FCont(promedioInv), rot),
	})

	// --- 11. DÍAS EN BODEGA ---
	dias := 0.0
	if rot > 0 {
		dias = 365 / rot
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Días de Inventario",
		Valor:          fmt.Sprintf("%.0f días", dias),
		Interpretacion: "Tiempo promedio del producto en bodega.",
		Clase:          "text-secondary",
		DetalleCuenta:  fmt.Sprintf("365 / Rotación (%.1f veces) = %.0f dias", rot, dias),
	})

	return dashboard, totales
}

// --- FUNCIONES DE APOYO ---
func evaluarIndice(valor, min, max float64) string {
	if valor < min {
		return "text-danger"
	}
	if valor <= max {
		return "text-success"
	}
	return "text-warning"
}

func evaluarIndiceEndeudamiento(v float64) string {
	if v > 70 {
		return "text-danger"
	}
	if v > 50 {
		return "text-warning"
	}
	return "text-success"
}

func GenerarDashboardIndustrial(costos ResumenCostos, ventas float64) []models.IndiceFinanciero {
	var dashboard []models.IndiceFinanciero

	// 1. COSTO PRIMO (Eficiencia de Inversión Directa)
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Costo Primo",
		Valor:          config.FCont(costos.CostoPrimo),
		Interpretacion: "Inversión directa en materiales y mano de obra.",
		Clase:          "text-primary",
		DetalleCuenta:  "MP Consumida + Mano de Obra Directa",
	})

	// 2. ABSORCIÓN DE COSTOS INDIRECTOS (CIF)
	// Idealmente el CIF no debe superar el 30-40% del costo total
	porcentajeCIF := 0.0
	if costos.CostoProduccion > 0 {
		porcentajeCIF = (costos.CIF / costos.CostoProduccion) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Impacto de Carga Fabril (CIF)",
		Valor:          fmt.Sprintf("%.2f%%", porcentajeCIF),
		Interpretacion: "Qué tanto pesan los gastos indirectos en la producción.",
		Clase:          evaluarIndice(porcentajeCIF, 0, 35), // Rojo si supera el 35%
		DetalleCuenta:  "(CIF / Costo Producción) * 100",
	})

	// 3. MARGEN INDUSTRIAL
	// ¿Cuánto nos queda después de fabricar antes de los gastos de admin?
	margenInd := 0.0
	if ventas > 0 {
		margenInd = ((ventas - costos.CostoProduccion) / ventas) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Margen Industrial",
		Valor:          fmt.Sprintf("%.2f%%", margenInd),
		Interpretacion: "Rentabilidad de la planta antes de gastos operativos.",
		Clase:          "text-success",
		DetalleCuenta:  "(Ventas - Costo Prod) / Ventas",
	})

	return dashboard
}
