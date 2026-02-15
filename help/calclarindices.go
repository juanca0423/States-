package help

import (
	"fmt"

	"ef/config"
	"ef/models"
)

func GenerarDashboard(activoCorr, pasivoCorr, inventario, costoVentas, activoTotal, pasivoTotal float64) []models.IndiceFinanciero {
	var dashboard []models.IndiceFinanciero

	// 1. LIQUIDEZ CORRIENTE
	liq := 0.0
	if pasivoCorr > 0 {
		liq = activoCorr / pasivoCorr
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Liquidez Corriente",
		Valor:          fmt.Sprintf("%.2f", liq),
		Interpretacion: "Capacidad de pago a corto plazo.",
		Clase:          evaluarIndice(liq, 1.0, 2.0),
		DetalleCuenta: fmt.Sprintf("Activo Corriente (%s) / Pasivo Corriente (%s)",
			config.FCont(activoCorr), config.FCont(pasivoCorr)),
	})

	// 2. ROTACIÓN DE INVENTARIOS
	rot := 0.0
	if inventario > 0 {
		rot = costoVentas / inventario
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Rotación de Inventario",
		Valor:          fmt.Sprintf("%.1f veces", rot),
		Interpretacion: "Veces que se renovó el stock en el periodo.",
		Clase:          "text-info",
		// Cambiamos %d por %s porque FCont devuelve un string
		DetalleCuenta: fmt.Sprintf("Costo de Ventas (%s) / Inventario (%s)", config.FCont(costoVentas), config.FCont(inventario)),
	})

	// 3. DÍAS DE INVENTARIO
	dias := 0.0
	if rot > 0 {
		dias = 360 / rot
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Días de Inventario",
		Valor:          fmt.Sprintf("%.0f días", dias),
		Interpretacion: "Tiempo promedio que el producto está en bodega.",
		Clase:          "text-secondary",
		// Usamos %f porque rot es float64
		DetalleCuenta: fmt.Sprintf("360 / Rotacion (%.2f)", rot),
	})

	// 4. ENDEUDAMIENTO TOTAL
	end := 0.0
	if activoTotal > 0 {
		end = (pasivoTotal / activoTotal) * 100
	}
	dashboard = append(dashboard, models.IndiceFinanciero{
		Nombre:         "Nivel de Endeudamiento",
		Valor:          fmt.Sprintf("%.2f%%", end),
		Interpretacion: "Porcentaje de activos financiados por terceros.",
		Clase:          evaluarIndiceEndeudamiento(end),
		DetalleCuenta:  fmt.Sprintf("(Pasivo %s / Activo %s) *100", config.FCont(pasivoTotal), config.FCont(activoTotal)),
	})

	return dashboard
}

// --- FUNCIONES DE APOYO ---

func evaluarIndice(valor, min, max float64) string {
	if valor < min {
		return "text-danger" // Insuficiente
	}
	if valor <= max {
		return "text-success" // Ideal
	}
	return "text-warning" // Exceso de liquidez (recursos ociosos)
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

// ... (Tus funciones GenerarDashboard, evaluarIndiceEndeudamiento, etc., están bien)

func CalcularTodo(costoVentas float64, invInicial float64, invFinal float64, actCorr float64, pasCorr float64) []models.IndiceFinanciero {
	var lista []models.IndiceFinanciero

	// 1. Calculamos la rotación
	promedioInv := (invInicial + invFinal) / 2
	rotacion := 0.0
	if promedioInv > 0 {
		rotacion = costoVentas / promedioInv
	}

	// 2. Calculamos los días que el producto pasa en bodega
	diasStock := 0.0
	if rotacion > 0 {
		diasStock = 365 / rotacion
	}

	// 3. Agregamos a la lista para usar las variables y que Go no dé error
	lista = append(lista, models.IndiceFinanciero{
		Nombre:         "Rotación de Inventario (KPI)",
		Valor:          fmt.Sprintf("%.2f veces", rotacion),
		Interpretacion: fmt.Sprintf("El inventario rota cada %.0f días", diasStock),
		Clase:          "text-info",
		DetalleCuenta: fmt.Sprintf("Activo Corriente (%s) / Pasivo Corriente (%s)",
			config.FCont(actCorr), config.FCont(pasCorr)),
	})

	// IMPORTANTE: El return que faltaba en la línea 162
	return lista
}

// Asegúrate de que todas las funciones tengan su return.
// CalcularIndices y GenerarDashboard ya lo tienen, así que con esto debería compilar.
