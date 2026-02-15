package help

import (
	"ef/config"
	"ef/models"
)

// HojaDeTrabajoIndustrial procesa las 9 columnas
func HojaDeTrabajoIndustrial(balanceRaw map[string]float64) []models.KV {
	var t models.Ht
	var listaKV []models.KV
	// Importante: Asegúrate de tener CreaMapCo en tu config
	catalogo := config.CreaMap(true)

	for codigoStr, valor := range balanceRaw {
		cuenta, existe := catalogo[codigoStr]
		if !existe || valor == 0 {
			continue
		}

		var ht models.HtString
		ht.Nombre = cuenta.Nombre

		switch cuenta.Saldo {
		case "activo":
			switch codigoStr {
			case "111301": // Mercaderías
				invFinal := balanceRaw["220005"]
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinal)
				ht.Perdidas, ht.Ganancias = config.FCont(valor), config.FCont(invFinal)
				t.Debe += valor
				t.Activo += invFinal
				t.Perdidas += valor
				t.Ganancias += invFinal

			case "111303", "111304": // MP y Proceso
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(valor)
				ht.CostoHaber = config.FCont(valor)
				t.Debe += valor
				t.Activo += valor
				t.CostoHaber += valor

			default:
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(valor)
				t.Debe += valor
				t.Activo += valor
			}

		case "pasivo":
			ht.Haber, ht.Pasivo = config.FCont(valor), config.FCont(valor)
			t.Haber += valor
			t.Pasivo += valor

		case "perdida":
			if codigoStr != "220001" {
				ht.Debe, ht.Perdidas = config.FCont(valor), config.FCont(valor)
				t.Debe += valor
				t.Perdidas += valor
			}

		case "ganancia":
			if codigoStr != "220005" {
				ht.Haber, ht.Ganancias = config.FCont(valor), config.FCont(valor)
				t.Haber += valor
				t.Ganancias += valor
			}

		case "costo_debe":
			ht.Debe, ht.CostoDebe = config.FCont(valor), config.FCont(valor)
			t.Debe += valor
			t.CostoDebe += valor

		case "costo_haber":
			ht.Haber, ht.CostoHaber = config.FCont(valor), config.FCont(valor)
			t.Haber += valor
			t.CostoHaber += valor
		}

		listaKV = append(listaKV, models.KV{Key: codigoStr, Value: ht})
	}

	return FinalizarHojaIndustrial(listaKV, t)
}

func FinalizarHojaIndustrial(lista []models.KV, t models.Ht) []models.KV {
	// Cálculo del Salto de Costo a Resultados
	costoProduccion := t.CostoDebe - t.CostoHaber

	filaCosto := models.KV{
		Key: "899999",
		Value: models.HtString{
			Nombre:     "COSTO DE PRODUCCIÓN",
			CostoHaber: config.FCont(costoProduccion),
			Perdidas:   config.FCont(costoProduccion),
		},
	}
	lista = append(lista, filaCosto)

	// Actualizamos totales para el cierre final
	t.Perdidas += costoProduccion

	return añadirTotalesYBalance(lista, t)
}
