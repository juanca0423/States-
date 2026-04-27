// Package config configuraciones de limpieza
package config

import (
	"sort"
	"strconv"

	"ef/models"
)

const (
	PorcentajeISR          = 0.25 // 25% sobre la utilidad
	PorcentajeReservaLegal = 0.05 // 5% sobre la utilidad neta de ISR
)

func LimpiaBalance(input map[string]any) map[string]float64 {
	output := make(map[string]float64)
	for k, v := range input {
		switch val := v.(type) {
		case float64:
			output[k] = val
		case string:
			f, _ := strconv.ParseFloat(val, 64)
			output[k] = f
		}
	}
	return output
}

func LimpiaMapa(formValues map[string][]string) map[string]float64 {
	balance := make(map[string]float64)
	for k, vv := range formValues {
		if len(vv) > 0 && vv[0] != "" {
			if f, err := strconv.ParseFloat(vv[0], 64); err == nil && f != 0 {
				balance[k] = f
			}
		}
	}
	return balance
}

func FinalizarHoja(lista []models.KV, t models.Ht) []models.KV {
	lista = append(lista, models.KV{Key: "900000", Value: models.HtString{
		Nombre: "SUMAS PRELIMINARES",
		Debe:   FCont(t.Debe), Haber: FCont(t.Haber),
		Perdidas: FCont(t.Perdidas), Ganancias: FCont(t.Ganancias),
		Activo: FCont(t.Activo), Pasivo: FCont(t.Pasivo),
		ClaTotal: "total",
	}})

	// Cálculo de Diferencia
	dif := t.Ganancias - t.Perdidas
	res := models.HtString{Nombre: "RESULTADO DEL EJERCICIO"}

	if dif > 0 { // Ganancia
		res.Perdidas, res.Pasivo = FCont(dif), FCont(dif)
		t.Perdidas += dif
		t.Pasivo += dif
	} else { // Pérdida
		res.Ganancias, res.Activo = FCont(-dif), FCont(-dif)
		t.Ganancias += (-dif)
		t.Activo += (-dif)
	}
	lista = append(lista, models.KV{Key: "910000", Value: res})

	lista = append(lista, models.KV{Key: "990000", Value: models.HtString{
		Nombre: "SUMAS IGUALES",
		Debe:   FCont(t.Debe), Haber: FCont(t.Haber),
		Perdidas: FCont(t.Perdidas), Ganancias: FCont(t.Ganancias),
		Activo: FCont(t.Activo), Pasivo: FCont(t.Pasivo),
		ClaTotal: "doble-linea",
	}})

	sort.Slice(lista, func(i, j int) bool { return lista[i].Key < lista[j].Key })
	return lista
}
