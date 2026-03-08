package models

// Estructura para cálculos (Matemática pura)

type BalanceNumerico struct {
	Nombre   string
	Valor1   float64 // Detalle
	Valor2   float64 // Sub-totales
	Valor3   float64 // NIIF (Corriente/No Corriente)
	Valor4   float64 // Gran Total
	EsCuenta bool    // Para saber si es cuenta o título
}

// Estructura para el HTML (Visual)

type BaStrin struct {
	Nombre     string
	Col1       string
	Col2       string
	Col3       string
	Col4       string
	ClasNombre string // Para indentación y negritas (clases CSS)
	ClasCol    string // Para colores de celdas o líneas
}

// Estructura pura para cálculos

type FilaBalance struct {
	Nombre         string
	Nivel          int     // 1: Título, 2: Rubro, 3: Cuenta, 4: Total
	V1, V2, V3, V4 float64 // Las 4 columnas de tu diseño
	EsNegativa     bool    // Para las depreciaciones/reservas
}

// Estructura para la vista (Handlebars)

type BalanceView struct {
	Nombre                 string
	Col1, Col2, Col3, Col4 string
	Clase                  string // Para el CSS personalizado
}

// Esta es la que usaremos para procesar tus variables

type Cuenta struct {
	Codigo string
	Nombre string
	Saldo  float64
}

// Esta es la que recibirá el HTML (.hbs)

type BaString struct {
	Codigo     string
	Nombre     string
	Col1       string
	Col2       string
	Col3       string
	Col4       string
	ClasNombre string // Para indentación y estilos
	Cla1       string
	Cla2       string
	Cla3       string
	Cla4       string
}

type TotalesBalance struct {
	ActivoCorriente   float64
	ActivoNoCorriente float64
	ActivoTotal       float64
	PasivoCorriente   float64
	PasivoNoCorriente float64
	PasivoTotal       float64
	Patrimonio        float64
	Inventario        float64 // Lo necesitamos para la Prueba Ácida y Rotación
	InventarioInicial float64 // <--- Agrega este para el Dashboard
}
