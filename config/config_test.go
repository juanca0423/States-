package config_test

import (
	"testing"

	"ef/config"
	"ef/models"
)

func setupMinimalNomenclatura() {
	config.Disponible = []models.Cue{
		{Codigo: 111101, Nombre: "Efectivo", Saldo: "activo", Categoria: "Disponible"},
	}
	config.Exigible = []models.Cue{
		{Codigo: 121001, Nombre: "Proveedores", Saldo: "pasivo", Categoria: "PasivoCorr"},
	}
	config.PatriNeto = []models.Cue{
		{Codigo: 131001, Nombre: "Capital", Saldo: "pasivo", Categoria: "PatriNeto"},
	}
	config.Ingresos = []models.Cue{
		{Codigo: 210001, Nombre: "Ventas", Saldo: "ganancia", Categoria: "Ingresos"},
	}
	config.Compras = []models.Cue{
		{Codigo: 220101, Nombre: "Compras", Saldo: "perdida", Categoria: "Compras", EsCosto: true},
	}
	config.CargarComercial()
}

func TestCargarComercialPopulatesMaps(t *testing.T) {
	setupMinimalNomenclatura()

	if len(config.ItemsCostos) == 0 {
		t.Fatal("expected ItemsCostos to be populated")
	}
	if len(config.ItemsContable) == 0 {
		t.Fatal("expected ItemsContable to be populated")
	}
	if len(config.MapCostos) == 0 {
		t.Fatal("expected MapCostos to be populated")
	}
	if len(config.MapContable) == 0 {
		t.Fatal("expected MapContable to be populated")
	}
}

func TestBuscarCuentaReturnsCuenta(t *testing.T) {
	setupMinimalNomenclatura()

	cuenta, ok := config.BuscarCuenta("111101", false)
	if !ok {
		t.Fatal("expected BuscarCuenta to return true for 111101")
	}
	if cuenta.Nombre != "Efectivo" {
		t.Fatalf("expected nombre Efectivo, got %q", cuenta.Nombre)
	}
}

func TestObtenerCuentasReturnsCostosWhenEsCostoTrue(t *testing.T) {
	setupMinimalNomenclatura()

	costos := config.ObtenerCuentas(true)
	contable := config.ObtenerCuentas(false)

	if len(costos) == 0 {
		t.Fatal("expected ObtenerCuentas(true) to return costo accounts")
	}
	if len(contable) == 0 {
		t.Fatal("expected ObtenerCuentas(false) to return contable accounts")
	}
	if len(costos) <= len(contable) {
		t.Fatal("expected cost accounts list to be non-empty and likely larger or distinct from contable list")
	}
}

func TestCreaMapReturnsProperMap(t *testing.T) {
	setupMinimalNomenclatura()

	m := config.CreaMap(true)
	if _, ok := m["111101"]; !ok {
		t.Fatal("expected CreaMap(true) to contain cost map entries")
	}
	m2 := config.CreaMap(false)
	if _, ok := m2["111101"]; !ok {
		t.Fatal("expected CreaMap(false) to contain contable map entries")
	}
}
