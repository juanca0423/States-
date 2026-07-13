package models

type IndiceFinanciero struct {
	Nombre         string
	Valor          string
	Interpretacion string // Para text-success, text-danger, etc.
	Clase          string
	DetalleCuenta  string // Guardará el texto de la fórmula
}

// DatosDashboard agrupa toda la información necesaria para los cálculos financieros
type DatosDashboard struct {
	ActivoCorriente   float64 // Disponibilidades, cuentas por cobrar, etc.
	PasivoCorriente   float64 // Deudas a corto plazo
	InventarioInicial float64 // Inventario al inicio del periodo
	InventarioFinal   float64 // Inventario al cierre del periodo
	CostoVentas       float64 // Costo de lo vendido
	ActivoTotal       float64 // Suma de todos los activos
	PasivoTotal       float64 // Suma de todas las deudas
	Ventas            float64 // Ingresos totales por ventas
	UtilidadNeta      float64 // Ganancia final después de impuestos
	GastosFijos       float64 // Costos que no varían con la producción
	GastosVariables   float64 // Costos que dependen del volumen
	GastosNoEfectivo  float64
}

type IndicesTotales struct {
	PuntoEContable     float64 // El punto de equilibrio normal
	PuntoECaja         float64 // El punto de equilibrio financiero (restando depreciaciones)
	MargenContribucion float64
	CostosTotales      float64
}
