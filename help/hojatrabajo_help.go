// Package help crea la hoja de trabajo
package help

import (
	"ef/config"
	"ef/models"
)

// Función auxiliar para no engordar la principal
func añadirTotalesYBalance(lista []models.KV, t models.Ht) []models.KV {
	// 1. Fila de Sumas Iguales
	suma := models.HtString{
		Nombre: "SUMAS IGUALES",
		Debe:   config.FCont(t.Debe), Haber: config.FCont(t.Haber),
		Perdidas: config.FCont(t.Perdidas), Ganancias: config.FCont(t.Ganancias),
		Activo: config.FCont(t.Activo), Pasivo: config.FCont(t.Pasivo),
		ClaTotal: "border-bottom-double",
	}
	lista = append(lista, models.KV{Key: "900000", Value: suma})

	// 2. Cálculo de Ganancia o Pérdida
	ganancia := t.Ganancias - t.Perdidas
	res := models.HtString{}

	if ganancia > 0 {
		res.Nombre = "GANANCIA ANTES DEL IMPUESTO"
		res.Perdidas = config.FCont(ganancia)
		res.Pasivo = config.FCont(ganancia)
		res.ClaTotal = "border-top"
	} else {
		res.Nombre = "PÉRDIDA DEL EJERCICIO"
		res.Ganancias = config.FCont(-ganancia)
		res.Activo = config.FCont(-ganancia)
		res.ClaTotal = "border-top"
	}

	return append(lista, models.KV{Key: "910000", Value: res})
}

// HojaDeTrabajo Cambiamos interface{} por any
func HojaDeTrabajo(balanceRaw map[string]float64) []models.KV {
	var t models.Ht
	var listaKV []models.KV

	catalogo := config.CreaMap(false)
	for codigoStr, valor := range balanceRaw {
		cuenta, existe := catalogo[codigoStr]
		if !existe {
			continue
		}
		var ht models.HtString
		ht.Nombre = cuenta.Nombre
		switch cuenta.Saldo {
		case "activo":
			if codigoStr == "111301" {
				invFinal := balanceRaw["220005"]
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinal)
				ht.Perdidas, ht.Ganancias = config.FCont(valor), config.FCont(invFinal)
				t.Debe += valor
				t.Activo += invFinal
				t.Perdidas += valor
				t.Ganancias += invFinal
			} else {
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(valor)
				t.Debe += valor
				t.Activo += valor
			}
		case "pasivo":
			ht.Haber, ht.Pasivo = config.FCont(valor), config.FCont(valor)
			t.Haber += valor
			t.Pasivo += valor
		case "perdida":
			if codigoStr == "220001" {
				continue
			}
			ht.Debe, ht.Perdidas = config.FCont(valor), config.FCont(valor)
			t.Debe += valor
			t.Perdidas += valor
		case "ganancia":
			if codigoStr == "220005" {
				continue
			}
			ht.Haber, ht.Ganancias = config.FCont(valor), config.FCont(valor)
			t.Haber += valor
			t.Ganancias += valor
		}
		listaKV = append(listaKV, models.KV{Key: codigoStr, Value: ht})
	}
	return config.FinalizarHoja(listaKV, t)
}
