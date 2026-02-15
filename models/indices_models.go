package models

type IndiceFinanciero struct {
	Nombre         string
	Valor          string
	Interpretacion string // Para text-success, text-danger, etc.
	Clase          string
	DetalleCuenta  string // Guardará el texto de la fórmula
}
