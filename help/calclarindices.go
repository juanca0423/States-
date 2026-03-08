package help

import (
	"fmt"

	"ef/config"
	"ef/models"
)

// GenerarDashboard centraliza todos los indicadores financieros
func GenerarDashboard(activoCorr, pasivoCorr, invInicial, invFinal, costoVentas, activoTotal, pasivoTotal, ventas, utilidadNeta float64) []models.IndiceFinanciero {
	var dashboard []models.IndiceFinanciero

	// --- 1. LIQUIDEZ CORRIENTE ---
	liq := 0.0
	if pasivoCorr > 0 {
		liq = activoCorr / pasivoCorr
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Liquidez Corriente",
		Valor:          fmt.Sprintf("%.2f", liq),
		Interpretacion: "Capacidad de cubrir deudas a corto plazo.",
		Clase:          evaluarIndice(liq, 1.0, 2.0),
		DetalleCuenta:  fmt.Sprintf("Act. Corriente (%s) / Pas. Corriente (%s)", config.FCont(activoCorr), config.FCont(pasivoCorr)),
	})

	// --- 2. PRUEBA DEL ÁCIDO (Nuevo) ---
	// Mide liquidez sin depender de la venta de inventario
	acido := 0.0
	if pasivoCorr > 0 {
		acido = (activoCorr - invFinal) / pasivoCorr
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Prueba del Ácido",
		Valor:          fmt.Sprintf("%.2f", acido),
		Interpretacion: "Liquidez inmediata (sin inventarios).",
		Clase:          evaluarIndice(acido, 0.8, 1.2),
		DetalleCuenta:  "(Act. Corr - Inv. Final) / Pas. Corr",
	})

	// --- 3. ROTACIÓN DE INVENTARIOS (Usando promedio) ---
	promedioInv := (invInicial + invFinal) / 2
	rot := 0.0
	if promedioInv > 0 {
		rot = costoVentas / promedioInv
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Rotación de Inventario",
		Valor:          fmt.Sprintf("%.1f veces", rot),
		Interpretacion: "Veces que el inventario se renovó en el año.",
		Clase:          "text-info",
		DetalleCuenta:  fmt.Sprintf("Costo Ventas (%s) / Inv. Promedio (%s)", config.FCont(costoVentas), config.FCont(promedioInv)),
	})

	// --- 4. DÍAS EN BODEGA ---
	dias := 0.0
	if rot > 0 {
		dias = 365 / rot
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Días de Inventario",
		Valor:          fmt.Sprintf("%.0f días", dias),
		Interpretacion: "Tiempo promedio del producto en bodega.",
		Clase:          "text-secondary",
		DetalleCuenta:  "365 / Rotación",
	})

	// --- 5. ENDEUDAMIENTO TOTAL ---
	end := 0.0
	if activoTotal > 0 {
		end = (pasivoTotal / activoTotal) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Nivel de Endeudamiento",
		Valor:          fmt.Sprintf("%.2f%%", end),
		Interpretacion: "Porcentaje de activos financiados por terceros.",
		Clase:          evaluarIndiceEndeudamiento(end),
		DetalleCuenta:  "(Pasivo Total / Activo Total) * 100",
	})

	// --- 6. MARGEN NETO (Nuevo) ---
	mNeto := 0.0
	if ventas > 0 {
		mNeto = (utilidadNeta / ventas) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Margen Neto de Utilidad",
		Valor:          fmt.Sprintf("%.2f%%", mNeto),
		Interpretacion: "Ganancia real por cada quetzal vendido.",
		Clase:          "text-success",
		DetalleCuenta:  "(Utilidad Neta / Ventas Totales) * 100",
	})

	return dashboard
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
