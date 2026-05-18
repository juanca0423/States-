// Package config config de cuenta
package config

import (
	"database/sql"
	"ef/models"
	"strconv"
)

var (
	Disponible         []models.Cue
	Exigible           []models.Cue
	Realisable         []models.Cue
	RealisableCo       []models.Cue
	PropPlanEqui       []models.Cue
	GtoIntan           []models.Cue
	GtoDiferidos       []models.Cue
	PasivoCorr         []models.Cue
	PasivoNoCorr       []models.Cue
	PatriNeto          []models.Cue
	Ingresos           []models.Cue
	InveIni            []models.Cue
	InveIniCo          []models.Cue
	Compras            []models.Cue
	InveFin            []models.Cue
	GtoVentas          []models.Cue
	GtoAdmin           []models.Cue
	IngrFina           []models.Cue
	GastosFina         []models.Cue
	OtrosGtoFina       []models.Cue
	Otring             []models.Cue
	MaterialesDirectos []models.Cue
	ManoObra           []models.Cue
	CostosIndirectos   []models.Cue
	CuentasInventarios []models.Cue
)

func CargarNomenclaturaDesdeDB(db *sql.DB) error {
	rows, err := db.Query("SELECT codigo, nombre, saldo, categoria, es_costo, es_variable, es_efectivo FROM nomenclatura ORDER BY codigo ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Limpiamos los slices por si acaso
	Disponible = nil
	Exigible = nil
	Realisable = nil
	RealisableCo = nil
	PropPlanEqui = nil
	GtoIntan = nil
	GtoDiferidos = nil
	PasivoCorr = nil
	PasivoNoCorr = nil
	PatriNeto = nil
	Ingresos = nil
	InveIni = nil
	InveIniCo = nil
	Compras = nil
	InveFin = nil
	GtoVentas = nil
	GtoAdmin = nil
	IngrFina = nil
	GastosFina = nil
	OtrosGtoFina = nil
	Otring = nil
	MaterialesDirectos = nil
	ManoObra = nil
	CostosIndirectos = nil
	CuentasInventarios = nil

	for rows.Next() {
		var c models.Cue

		err := rows.Scan(&c.Codigo, &c.Nombre, &c.Saldo, &c.Categoria, &c.EsCosto, &c.EsVariable, &c.EsEfectivo)
		if err != nil {
			return err
		}

		// El "Switch" mágico que llena tus variables actuales
		switch c.Categoria {
		case "Ingresos":
			Ingresos = append(Ingresos, c)
		case "InveIni":
			InveIni = append(InveIni, c)
		case "InveFin": // <--- ESTA ES LA QUE FALTA
			InveFin = append(InveFin, c)
		case "Compras":
			Compras = append(Compras, c)
		case "GtoVentas":
			GtoVentas = append(GtoVentas, c)
		case "GtoAdmin":
			GtoAdmin = append(GtoAdmin, c)
		case "IngrFina":
			IngrFina = append(IngrFina, c)
		case "GastosFina":
			GastosFina = append(GastosFina, c)
		case "OtrosGtoFina":
			OtrosGtoFina = append(OtrosGtoFina, c)
		case "Otring":
			Otring = append(Otring, c)
		case "Disponible":
			Disponible = append(Disponible, c)
		case "Exigible":
			Exigible = append(Exigible, c)
		case "Realisable":
			Realisable = append(Realisable, c)
		case "RealisableCo":
			RealisableCo = append(RealisableCo, c)
		case "PropPlanEqui":
			PropPlanEqui = append(PropPlanEqui, c)
		case "GtoIntan":
			GtoIntan = append(GtoIntan, c)
		case "GtoDiferidos":
			GtoDiferidos = append(GtoDiferidos, c)
		case "PasivoCorr":
			PasivoCorr = append(PasivoCorr, c)
		case "PasivoNoCorr":
			PasivoNoCorr = append(PasivoNoCorr, c)
		case "PatriNeto":
			PatriNeto = append(PatriNeto, c)
		// ... dentro del switch en CargarNomenclaturaDesdeDB ...
		case "MaterialesDirectos":
			MaterialesDirectos = append(MaterialesDirectos, c)
		case "ManoObra":
			ManoObra = append(ManoObra, c)
		case "CostosIndirectos":
			CostosIndirectos = append(CostosIndirectos, c)
		case "CuentasInventarios":
			CuentasInventarios = append(CuentasInventarios, c)
		}
	}
	// ¡ESTO ES LO NUEVO!
	// Una vez llenos los slices, regeneramos los mapas y listas de items
	MapContable = make(map[string]CuentaItem)
	MapCostos = make(map[string]CuentaItem)
	CargarComercial()
	CargarCostos()
	return nil
}

type CuentaItem struct {
	Codigo     string `json:"codigo"`
	CodInt     int    `json:"-"`
	Nombre     string `json:"nombre"`
	Saldo      string `json:"saldo"` // Asegúrate que sea String si en la DB es VARCHAR
	Categoria  string `json:"categoria"`
	EsCosto    bool   `json:"es_costo"`
	EsEfectivo bool   `json:"es_efectivo"`
	EsVariable bool   `json:"es_variable"`
}

var (
	ItemsContable []CuentaItem
	ItemsCostos   []CuentaItem
	MapContable   map[string]CuentaItem
	MapCostos     map[string]CuentaItem
)

func CargarComercial() {
	// 1. Limpiar todo
	ItemsContable = []CuentaItem{}
	ItemsCostos = []CuentaItem{}
	MapContable = make(map[string]CuentaItem)
	MapCostos = make(map[string]CuentaItem)

	// 2. Todos los grupos que manejas en el Switch
	todosLosGrupos := [][]models.Cue{
		Disponible, Exigible, Realisable, PropPlanEqui,
		GtoIntan, GtoDiferidos, PasivoCorr, PasivoNoCorr,
		PatriNeto, Ingresos, InveIni, Compras,
		InveFin, GtoVentas, GtoAdmin, IngrFina,
		GastosFina, OtrosGtoFina, Otring,
		MaterialesDirectos, ManoObra, CostosIndirectos, CuentasInventarios,
	}

	for _, grupo := range todosLosGrupos {
		for _, v := range grupo {
			key := strconv.Itoa(v.Codigo)

			item := CuentaItem{
				Codigo:     key,
				CodInt:     v.Codigo,
				Nombre:     v.Nombre,
				Saldo:      v.Saldo,
				Categoria:  v.Categoria,
				EsCosto:    v.EsCosto,
				EsEfectivo: v.EsEfectivo,
				EsVariable: v.EsVariable,
			}

			// --- LÓGICA PARA REPORTES (Comercial) ---
			// Solo entra si NO es costo
			if !v.EsCosto {
				if _, existe := MapContable[key]; !existe {
					ItemsContable = append(ItemsContable, item)
					MapContable[key] = item
				}
			}

			// --- LÓGICA PARA ADMIN (Panel de Gestión) ---
			// Entran ABSOLUTAMENTE TODAS (las 66)
			if _, existe := MapCostos[key]; !existe {
				ItemsCostos = append(ItemsCostos, item)
				MapCostos[key] = item
			}
		}
	}
}

func CargarCostos() {
	// Ya no necesita lógica propia porque CargarComercial llena MapCostos con todo.
}

func ObtenerCuentas(esCosto bool) []CuentaItem {
	if esCosto {
		return ItemsCostos
	}
	return ItemsContable
}

func BuscarCuenta(codigo string, esCosto bool) (CuentaItem, bool) {
	var m map[string]CuentaItem
	if esCosto {
		m = MapCostos
	} else {
		m = MapContable
	}
	cuenta, ok := m[codigo]
	return cuenta, ok
}

var (
	DBDatos    map[string]models.Cue
	DBDatosCos map[string]models.Cue
)

func CreaMap(esCosto bool) map[string]CuentaItem {
	if esCosto {
		return MapCostos
	}
	return MapContable
}
