package config

import (
	"math"

	"ef/models"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FCont(v float64) string {
	if v == 0 {
		return "-" // O "-" si prefieres el estilo contable para valores cero
	}

	negativo := v < 0
	v = math.Abs(v)

	// Formatear con comas y decimales
	p := message.NewPrinter(language.English)
	resultado := p.Sprintf("%.2f", v)

	if negativo {
		return "(" + resultado + ")"
	}
	return resultado
}

func DividirCuentas(out map[string]float64) map[string]models.Cuenta {
	balance := make(map[string]models.Cuenta)

	// b es tu nomenclatura (CreaMap)
	b := CreaMap(false)

	for k, v := range out {
		// Buscamos el nombre en la nomenclatura
		nombre := "Cuenta no definida"
		if info, existe := b[k]; existe {
			nombre = info.Nombre
		}
		// Creamos el detalle
		detalle := models.Cuenta{
			Codigo: k,
			Nombre: nombre,
			Saldo:  v,
		}
		balance[k] = detalle
	}
	return balance
}
