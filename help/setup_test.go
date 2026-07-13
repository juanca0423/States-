package help

import (
	"ef/config"
	"ef/models"
)

func SetupNomenclatura() {
	// --- ACTIVOS ---
	config.Disponible = []models.Cue{
		{Codigo: 111101, Nombre: "Efectivo", Saldo: "activo", Categoria: "Disponible"},
		{Codigo: 111102, Nombre: "Bancos", Saldo: "activo", Categoria: "Disponible"},
	}
	config.Exigible = []models.Cue{
		{Codigo: 111201, Nombre: "Deu. no Comerciales", Saldo: "activo", Categoria: "Exigible"},
		{Codigo: 111203, Nombre: "Clientes", Saldo: "activo", Categoria: "Exigible"},
		{Codigo: 111204, Nombre: "(-)Res. para Cuen Inco.", Saldo: "pasivo", Categoria: "Exigible"},
	}
	config.Realisable = []models.Cue{
		{Codigo: 111301, Nombre: "Mercaderias", Saldo: "activo", Categoria: "Realisable"},
		{Codigo: 111303, Nombre: "Materia Prima", Saldo: "activo", Categoria: "Realisable", EsCosto: true},
		{Codigo: 111304, Nombre: "Productos en Proceso", Saldo: "activo", Categoria: "Realisable", EsCosto: true},
		{Codigo: 111305, Nombre: "Productos Terminados", Saldo: "activo", Categoria: "Realisable", EsCosto: true},
	}
	config.PropPlanEqui = []models.Cue{
		{Codigo: 112101, Nombre: "Mobiliario y Equipo", Saldo: "activo", Categoria: "PropPlanEqui"},
		{Codigo: 112102, Nombre: "(-) Dep. Acu. Mob y Equipo", Saldo: "pasivo", Categoria: "PropPlanEqui"},
	}

	// --- PASIVOS Y CAPITAL ---
	config.PasivoCorr = []models.Cue{
		{Codigo: 121001, Nombre: "Proveedores", Saldo: "pasivo", Categoria: "PasivoCorr"},
		{Codigo: 121005, Nombre: "IGSS por Pagar", Saldo: "pasivo", Categoria: "PasivoCorr"},
	}
	config.PatriNeto = []models.Cue{
		{Codigo: 131001, Nombre: "Capital", Saldo: "pasivo", Categoria: "PatriNeto"},
		{Codigo: 131003, Nombre: "(-)Cuenta Personal", Saldo: "activo", Categoria: "PatriNeto"},
	}

	// --- RESULTADOS (INGRESOS Y GASTOS) ---
	config.Ingresos = []models.Cue{
		{Codigo: 210001, Nombre: "Ventas", Saldo: "ganancia", Categoria: "Ingresos"},
	}
	config.InveIni = []models.Cue{
		{Codigo: 220001, Nombre: "Inve. Ini de Mercaderías", Saldo: "perdida", Categoria: "InveIni"},
	}
	config.Compras = []models.Cue{
		{Codigo: 220101, Nombre: "Compras", Saldo: "perdida", Categoria: "Compras"},
	}

	// --- GASTOS DE OPERACIÓN ---
	config.GtoVentas = []models.Cue{
		{Codigo: 231001, Nombre: "Sueldos Sala de Ventas", Saldo: "perdida", Categoria: "GtoVentas"},
		{Codigo: 231004, Nombre: "Dep. Mob. y Equi S/Ventas", Saldo: "perdida", Categoria: "GtoVentas"},
	}
	config.GtoAdmin = []models.Cue{
		{Codigo: 232001, Nombre: "Sueldos de Admón.", Saldo: "perdida", Categoria: "GtoAdmin"},
		{Codigo: 232008, Nombre: "Cuentas Incobrables", Saldo: "perdida", Categoria: "GtoAdmin"},
	}

	// --- COSTOS DE PRODUCCIÓN (Materia Prima y Mano de Obra) ---
	// Usamos saldo "costo_debe" según tu nomenclatura
	config.MaterialesDirectos = []models.Cue{
		{Codigo: 310100, Nombre: "Inventario Inicial MP", Saldo: "costo_debe", Categoria: "MaterialesDirectos", EsCosto: true},
		{Codigo: 310200, Nombre: "Compras MP", Saldo: "costo_debe", Categoria: "MaterialesDirectos", EsCosto: true},
		{Codigo: 310800, Nombre: "(-) Inv. Final MP", Saldo: "costo_haber", Categoria: "MaterialesDirectos", EsCosto: true},
	}
	config.ManoObra = []models.Cue{
		{Codigo: 320100, Nombre: "Sueldos y Salarios Fábrica", Saldo: "costo_debe", Categoria: "ManoObra", EsCosto: true},
	}
	config.CuentasInventarios = []models.Cue{
		{Codigo: 340101, Nombre: "Inv. Inicial Proceso", Saldo: "costo_debe", Categoria: "CuentasInventarios", EsCosto: true},
		{Codigo: 340201, Nombre: "(-) Inv. Final Proceso", Saldo: "costo_haber", Categoria: "CuentasInventarios", EsCosto: true},
	}

	// Agregamos Inventario en Tránsito
	config.Realisable = append(config.Realisable, models.Cue{
		Codigo:    111306,
		Nombre:    "Inventarios en Tránsito",
		Saldo:     "activo",
		Categoria: "Realisable",
		EsCosto:   false,
	})

	config.CargarComercial()
}
