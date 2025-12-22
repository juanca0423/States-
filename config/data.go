package config

import (
	"ef/models"
	"strconv"
)

func CreaCuentas() []models.Cue {

	Disponible := []models.Cue{
		{
			Codigo: 111101,
			Nombre: "Efectivo",
			Saldo:  "activo",
			Clases: "",
			Valor:  6510},
		{
			Codigo: 111102,
			Nombre: "Bancos",
			Saldo:  "activo",
			Clases: "",
			Valor:  197650},
	}

	Exigible := []models.Cue{
		{
			Codigo: 111201,
			Nombre: "Deu. no Comerciales",
			Saldo:  "activo",
			Clases: "",
			Valor:  3000,
		},
		{
			Codigo: 111202,
			Nombre: "Doc. por Cobrar",
			Saldo:  "activo",
			Clases: "",
			Valor:  0,
		},
		{
			Codigo: 111203,
			Nombre: "Clientes",
			Saldo:  "activo",
			Clases: "",
			Valor:  40700},
		{
			Codigo: 111204,
			Nombre: "(-)Res. para Cuen Inco.",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  1221,
		},
	}

	Realisable := []models.Cue{
		{
			Codigo: 111301,
			Nombre: "Mercaderías",
			Saldo:  "activo",
			Clases: "",
			Valor:  44000,
		},
		{
			Codigo: 111302,
			Nombre: "Papelería y útiles",
			Saldo:  "activo",
			Clases: "",
			Valor:  1320,
		},
	}
	PropPlanEqui := []models.Cue{
		{
			Codigo: 112101,
			Nombre: "Mobiliario y Equipo",
			Saldo:  "activo",
			Clases: "",
			Valor:  10750,
		},
		{
			Codigo: 112102,
			Nombre: "(-) Dep. Acu. Mob y Equipo",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  2050,
		},
		{
			Codigo: 112103,
			Nombre: "Equipo de Computo",
			Saldo:  "activo",
			Clases: "",
			Valor:  28700,
		},
		{
			Codigo: 112104,
			Nombre: "Dep. Acu. de Equi de Compu",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  9232.41,
		},
	}
	PasivoCorr := []models.Cue{
		{
			Codigo: 121001,
			Nombre: "Proveedores",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  38800,
		},
		{
			Codigo: 121002,
			Nombre: "IVA",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  5140,
		},
		{
			Codigo: 121003,
			Nombre: "Doc. por Pagar a C/Plazo",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  14000,
		},
		{
			Codigo: 121004,
			Nombre: "Préstamos Banc C/Plazo",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  80000,
		},
		{
			Codigo: 121005,
			Nombre: "IGSS por Pagar",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  2117.5,
		},
		{
			Codigo: 121006,
			Nombre: "Comi. Percibi no Deven",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  2040,
		},
		{
			Codigo: 121007,
			Nombre: "ISR por Pagar",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  0,
		},
		{
			Codigo: 121008,
			Nombre: "(-) ISR Trimestral",
			Saldo:  "activo",
			Clases: "",
			Valor:  10800,
		},
	}

	GtoIntan := []models.Cue{
		{
			Codigo: 112201,
			Nombre: "Gastos de Organización",
			Saldo:  "activo",
			Clases: "",
			Valor:  752,
		},
		{
			Codigo: 112202,
			Nombre: "Amor. Acu. de Gtos de Organ",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  150.4,
		},
	}

	GtoDiferidos := []models.Cue{
		{
			Codigo: 112301,
			Nombre: "Inve. en Valores a L P",
			Saldo:  "activo",
			Clases: "",
			Valor:  100000,
		},
	}
	PasivoNoCorr := []models.Cue{
		{
			Codigo: 122001,
			Nombre: "Doc. por Pagar a L P",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  10000,
		},
	}

	PatriNeto := []models.Cue{
		{
			Codigo: 131001,
			Nombre: "Capital",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  109409.8,
		},
		{
			Codigo: 131002,
			Nombre: "Ganancia después del ISR",
			Saldo:  "pasivo",
			Clases: "",
			Valor:  0,
		},
		{
			Codigo: 131003,
			Nombre: "(-)Cuenta Personal",
			Saldo:  "activo",
			Clases: "",
			Valor:  60000,
		},
	}
	Ingresos := []models.Cue{
		{
			Codigo: 210001,
			Nombre: "Ventas Brutas",
			Saldo:  "ganancia",
			Clases: "",
			Valor:  540434.54,
		},
		{
			Codigo: 210002,
			Nombre: "(-)Dev y Rebajas S/Ventas",
			Saldo:  "perdida",
			Clases: "",
			Valor:  3750,
		},
	}
	CostoVentas := []models.Cue{
		{
			Codigo: 220001,
			Nombre: "Inve. Ini de Mercaderías",
			Saldo:  "perdida",
			Clases: "",
			Valor:  44000,
		},
		{
			Codigo: 220002,
			Nombre: "Compras",
			Saldo:  "perdida",
			Clases: "",
			Valor:  95800,
		},
		{
			Codigo: 220003,
			Nombre: "Gastos sobre Compras",
			Saldo:  "perdida",
			Clases: "",
			Valor:  5640,
		},
		{
			Codigo: 220004,
			Nombre: "Dev. y Reba. sobre Compras",
			Saldo:  "ganancia",
			Clases: "",
			Valor:  2740,
		},
		{
			Codigo: 220005,
			Nombre: "Inve. Final de Mercaderías",
			Saldo:  "ganancia",
			Clases: "",
			Valor:  63740,
		},
	}

	GtoVentas := []models.Cue{
		{
			Codigo: 231001,
			Nombre: "Sueldos Sala de Ventas",
			Saldo:  "perdida",
			Clases: "",
			Valor:  70800,
		},
		{
			Codigo: 231002,
			Nombre: "Boni Incen Sala de Ventas",
			Saldo:  "perdida",
			Clases: "",
			Valor:  9000,
		},
		{
			Codigo: 231003,
			Nombre: "Cuo. Patro S/Ventas",
			Saldo:  "perdida",
			Clases: "",
			Valor:  8970.36,
		},
		{
			Codigo: 231004,
			Nombre: "Dep. Mob. y Equi S/Ventas",
			Saldo:  "perdida",
			Clases: "",
			Valor:  1230,
		},
	}

	GtoAdmin := []models.Cue{
		{
			Codigo: 232001,
			Nombre: "Sueldos de Admón.",
			Saldo:  "perdida",
			Clases: "",
			Valor:  74400,
		},
		{
			Codigo: 232002,
			Nombre: "Boni Incentivo de Admon",
			Saldo:  "perdida",
			Clases: "",
			Valor:  12000,
		},
		{
			Codigo: 232003,
			Nombre: "Cuo Pat de Admon",
			Saldo:  "perdida",
			Clases: "",
			Valor:  9426.48,
		},
		{
			Codigo: 232004,
			Nombre: "Dep. Mob y Equi de Oficina",
			Saldo:  "perdida",
			Clases: "",
			Valor:  820,
		},
		{
			Codigo: 232005,
			Nombre: "Dep Equi de Computación",
			Saldo:  "perdida",
			Clases: "",
			Valor:  9232.41,
		},
		{
			Codigo: 232006,
			Nombre: "Amor. Gtos de Organización",
			Saldo:  "perdida",
			Clases: "",
			Valor:  155.4,
		},
		{
			Codigo: 232007,
			Nombre: "Pap y Utiles Consumidos",
			Saldo:  "perdida",
			Clases: "",
			Valor:  1838,
		},
		{
			Codigo: 232008,
			Nombre: "Cuentas Incobrables",
			Saldo:  "perdida",
			Clases: "",
			Valor:  1221,
		},
	}

	IngrFina := []models.Cue{
		{
			Codigo: 240001,
			Nombre: "Intereses Percibidos",
			Saldo:  "ganancia",
			Clases: "",
			Valor:  9780,
		},
	}

	GastosFina := []models.Cue{
		{
			Codigo: 250001,
			Nombre: "Intereses Gastos",
			Saldo:  "perdida",
			Clases: "",
			Valor:  8650,
		},
	}

	OtrosGtoFina := []models.Cue{
		{
			Codigo: 260001,
			Nombre: "Donativos",
			Saldo:  "perdida",
			Clases: "",
			Valor:  10000,
		},
	}

	var Cuentas []models.Cue

	Cuentas = append(Cuentas, Disponible...)
	Cuentas = append(Cuentas, Exigible...)
	Cuentas = append(Cuentas, Realisable...)
	Cuentas = append(Cuentas, PropPlanEqui...)
	Cuentas = append(Cuentas, GtoIntan...)
	Cuentas = append(Cuentas, GtoDiferidos...)
	Cuentas = append(Cuentas, PasivoCorr...)
	Cuentas = append(Cuentas, PasivoNoCorr...)
	Cuentas = append(Cuentas, PatriNeto...)
	Cuentas = append(Cuentas, Ingresos...)
	Cuentas = append(Cuentas, CostoVentas...)
	Cuentas = append(Cuentas, GtoVentas...)
	Cuentas = append(Cuentas, GtoAdmin...)
	Cuentas = append(Cuentas, IngrFina...)
	Cuentas = append(Cuentas, GastosFina...)
	Cuentas = append(Cuentas, OtrosGtoFina...)
	return Cuentas
}

func CreaMap() map[string]models.Cue {
	var DbDatos = make(map[string]models.Cue)
	Cuentas := CreaCuentas()

	for _, v := range Cuentas {
		key := strconv.Itoa(v.Codigo)
		DbDatos[key] = v
	}
	return DbDatos
}
