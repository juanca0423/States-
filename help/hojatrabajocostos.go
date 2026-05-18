package help

import (
	"ef/config"
	"ef/models"
)

// HojaDeTrabajoIndustrial procesa las 9 columnas
func HojaDeTrabajoIndustrial(balanceRaw map[string]float64) []models.KV {
	var t models.Ht
	var listaKV []models.KV
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
				invFinalMerc := balanceRaw["220005"]
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinalMerc)
				ht.Perdidas, ht.Ganancias = config.FCont(valor), config.FCont(invFinalMerc)
				t.Debe += valor
				t.Activo += invFinalMerc
				t.Perdidas += valor
				t.Ganancias += invFinalMerc
			case "111303": // MP
				invFinalMP := balanceRaw["310800"]
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinalMP)
				ht.CostoHaber, ht.CostoDebe = config.FCont(invFinalMP), config.FCont(valor)
				t.Debe += valor
				t.Activo += invFinalMP
				t.CostoHaber += valor
				t.CostoDebe += invFinalMP
			case "111304": // Proceso
				invFinalProceso := balanceRaw["340201"]
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinalProceso)
				ht.CostoDebe, ht.CostoHaber = config.FCont(invFinalProceso), config.FCont(valor)
				t.Debe += valor
				t.Activo += invFinalProceso
				t.CostoHaber += valor
				t.CostoDebe += invFinalProceso
			case "111305": // PRODUCTOS TERMINADOS (La clave)
				invFinalPT := balanceRaw["340202"] // Nueva cuenta
				ht.Debe, ht.Activo = config.FCont(valor), config.FCont(invFinalPT)
				// VAN A RESULTADOS, NO A COSTO
				ht.Perdidas, ht.Ganancias = config.FCont(valor), config.FCont(invFinalPT)
				t.Perdidas += valor
				t.Ganancias += invFinalPT
				t.Activo += invFinalPT
				t.Debe += valor
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
