package help

import (
	"testing"

	"ef/config"
	"ef/models"
)

// grupoPrueba construye un catálogo pequeño y un Balance para ejercitar RecCol*.
func grupoPrueba() ([]models.Cue, map[string]float64) {
	cuentas := []models.Cue{
		{Codigo: 231001, Nombre: "Sueldos Sala de Ventas", Saldo: "perdida", EsVariable: true, EsEfectivo: true},
		{Codigo: 231004, Nombre: "Dep. Mob. y Equi S/Ventas", Saldo: "perdida", EsVariable: false, EsEfectivo: false},
	}
	balance := map[string]float64{
		"231001": 1000.0,
		"231004": 500.0,
	}
	return cuentas, balance
}

// últimaFila devuelve el Value de la última fila (la del total del grupo).
func últimaFila(res []models.KR) models.ReString {
	return res[len(res)-1].Value
}

func TestRecCol1SinCol3(t *testing.T) {
	cuentas, balance := grupoPrueba()

	var vVar, vFij, vNoEf float64
	res, saldo := RecCol1(balance, cuentas, "Gastos de Distribución", &vVar, &vFij, &vNoEf)

	if saldo != 1500 {
		t.Fatalf("SaldoGrupo esperado 1500, obtuvo %.2f", saldo)
	}
	// título + 2 filas de datos
	if len(res) != 3 {
		t.Fatalf("filas esperadas 3, obtuvo %d", len(res))
	}
	if res[0].Key != "TIT-Gastos de Distribución" {
		t.Errorf("clave de título incorrecta: %q", res[0].Key)
	}

	ult := últimaFila(res)
	if ult.Col2 != config.FCont(1500) {
		t.Errorf("Col2 incorrecta: %q", ult.Col2)
	}
	if ult.Cla1 != "border-bottom" {
		t.Errorf("Cla1 incorrecta: %q", ult.Cla1)
	}
	if ult.Col3 != "" {
		t.Errorf("RecCol1 no debe pintar Col3, obtuvo %q", ult.Col3)
	}

	// acumuladores de costos
	if vVar != 1000 {
		t.Errorf("Variables esperado 1000, obtuvo %.2f", vVar)
	}
	if vFij != 500 {
		t.Errorf("Fijos esperado 500, obtuvo %.2f", vFij)
	}
	if vNoEf != 500 { // solo 231004 es no-efectivo
		t.Errorf("NoEfectivo esperado 500, obtuvo %.2f", vNoEf)
	}
}

func TestRecCol1TotConSuma(t *testing.T) {
	cuentas, balance := grupoPrueba()

	var vVar, vFij, vNoEf float64
	res, saldo := RecCol1Tot(balance, cuentas, "Gastos de Administración", 2000, &vVar, &vFij, &vNoEf)

	if saldo != 1500 {
		t.Fatalf("SaldoGrupo esperado 1500, obtuvo %.2f", saldo)
	}
	ult := últimaFila(res)
	if ult.Col2 != config.FCont(1500) {
		t.Errorf("Col2 incorrecta: %q", ult.Col2)
	}
	if ult.Col3 != config.FCont(3500) { // 2000 + 1500
		t.Errorf("Col3 incorrecta: %q", ult.Col3)
	}
	if ult.Cla3 != "border-bottom" || ult.Cla2 != "border-bottom" || ult.Cla1 != "border-bottom" {
		t.Errorf("bordes incorrectos: Cla1=%q Cla2=%q Cla3=%q", ult.Cla1, ult.Cla2, ult.Cla3)
	}
}

func TestRecCol2TotConResta(t *testing.T) {
	cuentas, balance := grupoPrueba()

	var vVar, vFij, vNoEf float64
	res, saldo := RecCol2Tot(balance, cuentas, "Gastos Financieros", 2000, &vVar, &vFij, &vNoEf)

	if saldo != 1500 {
		t.Fatalf("SaldoGrupo esperado 1500, obtuvo %.2f", saldo)
	}
	ult := últimaFila(res)
	if ult.Col3 != config.FCont(500) { // 2000 - 1500
		t.Errorf("Col3 incorrecta: %q", ult.Col3)
	}
	if ult.Cla3 != "border-bottom" {
		t.Errorf("Cla3 incorrecta: %q", ult.Cla3)
	}
}

// TestRecColGrupoEquivaleAWrappers garantiza que el núcleo unificado produce
// exactamente la misma salida que el wrapper público (regresión del refactor).
func TestRecColGrupoEquivaleAWrappers(t *testing.T) {
	cuentas, balance := grupoPrueba()

	var a1, b1, c1 float64
	r1, s1 := RecCol1(balance, cuentas, "G", &a1, &b1, &c1)

	var a2, b2, c2 float64
	r2, s2 := recColGrupo(balance, cuentas, "G", 0, nil, &a2, &b2, &c2)

	if s1 != s2 {
		t.Fatalf("Saldos distintos: %f vs %f", s1, s2)
	}
	if len(r1) != len(r2) {
		t.Fatalf("largos distintos: %d vs %d", len(r1), len(r2))
	}
	for i := range r1 {
		if r1[i].Key != r2[i].Key || r1[i].Value != r2[i].Value {
			t.Errorf("fila %d distinta: %+v vs %+v", i, r1[i], r2[i])
		}
	}
}

func TestRecColIgnoraCeros(t *testing.T) {
	cuentas, _ := grupoPrueba()
	balance := map[string]float64{"231001": 0.0, "231004": 0.0}

	var vVar, vFij, vNoEf float64
	res, saldo := RecCol1(balance, cuentas, "G", &vVar, &vFij, &vNoEf)

	if res != nil || saldo != 0 {
		t.Errorf("esperado nil y saldo 0, obtuvo res=%v saldo=%.2f", res, saldo)
	}
}
