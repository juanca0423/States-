package help

import (
	"testing"

	"ef/config"
	"ef/models"
)

func TestCuadreBalance(t *testing.T) {
	// 1. LIMPIEZA PREVENTIVA: Aseguramos que el config esté como queremos
	config.Ingresos = []models.Cue{{Codigo: 210001, Saldo: "ganancia"}}
	config.InveIni = []models.Cue{{Codigo: 220001, Saldo: "perdida"}}
	config.Compras = []models.Cue{{Codigo: 220101, Saldo: "perdida"}}

	// 1. MOCK del Config (Igual que hicimos con resultados)
	config.Disponible = []models.Cue{{Codigo: 111101, Nombre: "Efectivo", Saldo: "activo"}}
	config.Exigible = []models.Cue{{Codigo: 111203, Nombre: "Clientes", Saldo: "activo"}, {Codigo: 111204, Nombre: "(-)Reserva", Saldo: "pasivo"}}
	config.PasivoCorr = []models.Cue{{Codigo: 121001, Nombre: "Proveedores", Saldo: "pasivo"}}
	config.PatriNeto = []models.Cue{{Codigo: 131001, Nombre: "Capital", Saldo: "pasivo"}}

	// 2. Datos de prueba:
	// Activo: Efectivo(500) + Clientes(200) - Reserva(50) = 650
	// Pasivo: Proveedores(300)
	// Utilidad Neta (que viene de resultados): 100
	// Patrimonio Necesario para cuadrar: 650 - 300 - 100 = 250
	saldos := map[string]models.Cuenta{
		"111101": {Saldo: 500},
		"111203": {Saldo: 200},
		"111204": {Saldo: 50}, // Esta debe restar al activo por ser saldo "pasivo"
		"121001": {Saldo: 300},
		"131001": {Saldo: 250},
	}
	utilidadNeta := 100.0

	// 3. Ejecutamos la lógica
	_, totales := GenerarTodoElBalance(saldos, utilidadNeta)

	// 4. Verificaciones
	// El Activo debe ser 650
	esperadoActivo := 650.0
	if totales.ActivoTotal != esperadoActivo {
		t.Errorf("ACTIVO INCORRECTO: Se obtuvo %.2f, se esperaba %.2f", totales.ActivoTotal, esperadoActivo)
	}

	// El Pasivo + Patrimonio + Utilidad debe sumar igual al Activo
	sumaLadoDerecho := totales.PasivoTotal + totales.Patrimonio + utilidadNeta
	if sumaLadoDerecho != totales.ActivoTotal {
		t.Errorf("BALANCE DESCUADRADO: Activo(%.2f) != Pasivo+Pat(%.2f)", totales.ActivoTotal, sumaLadoDerecho)
	}
}

func TestResultados(t *testing.T) {
	// 1. LIMPIEZA PREVENTIVA: Aseguramos que el config esté como queremos
	config.Ingresos = []models.Cue{{Codigo: 210001, Saldo: "ganancia"}}
	config.InveIni = []models.Cue{{Codigo: 220001, Saldo: "perdida"}}
	config.Compras = []models.Cue{{Codigo: 220101, Saldo: "perdida"}}

	// ... resto de tu lógica de balancePrueba y Resultados()
	// 1. IMPORTANTE: Rellenar el config para que las funciones encuentren los códigos
	config.Ingresos = []models.Cue{{Codigo: 210001, Nombre: "Ventas", Saldo: "ganancia"}}
	config.InveIni = []models.Cue{{Codigo: 220001, Nombre: "Inv. Inicial", Saldo: "perdida"}}
	config.Compras = []models.Cue{{Codigo: 220101, Nombre: "Compras", Saldo: "perdida"}}

	// 2. Datos del test
	balancePrueba := map[string]float64{
		"210001": 1000.00, // Ventas
		"220001": 200.00,  // Inv. Inicial
		"220101": 300.00,  // Compras
		"220005": 100.00,  // Inv. Final (este se busca directo, no necesita config)
	}

	// 3. Ejecutar
	_, totales := Resultados(balancePrueba, 0)

	// 4. Verificaciones
	if totales.VentasNetas != 1000 {
		t.Errorf("FALLO VENTAS: Se obtuvo %.2f, se esperaba 1000.00", totales.VentasNetas)
	}

	// Costo esperado: (200 + 300) - 100 = 400
	if totales.CostoVentas != 400 {
		t.Errorf("FALLO COSTO: Se obtuvo %.2f, se esperaba 400.00", totales.CostoVentas)
	}

	// Utilidad Bruta: 1000 - 400 = 600
	if (totales.VentasNetas - totales.CostoVentas) != 600 {
		t.Errorf("FALLO UTILIDAD: Se obtuvo %.2f", totales.VentasNetas-totales.CostoVentas)
	}
}

func TestGenerarDashboard(t *testing.T) {
	t.Run("Escenario Ideal", func(t *testing.T) {
		// Datos: ActCorr=200, PasCorr=100, InvIni=50, InvFin=50, Costo=400, ActTotal=1000, PasTotal=400, Ventas=1000, Utilidad=100
		indices := GenerarDashboard(200, 100, 50, 50, 400, 1000, 400, 1000, 100)

		// 1. Verificar Liquidez (200/100 = 2.00)
		if indices[0].Valor != "2.00" {
			t.Errorf("Liquidez errónea: se obtuvo %s, se esperaba 2.00", indices[0].Valor)
		}

		// 2. Verificar Prueba del Ácido ((200-50)/100 = 1.50)
		if indices[1].Valor != "1.50" {
			t.Errorf("Ácido erróneo: se obtuvo %s, se esperaba 1.50", indices[1].Valor)
		}

		// 3. Verificar Rotación (Costo 400 / Promedio 50 = 8.0)
		if indices[2].Valor != "8.0 veces" {
			t.Errorf("Rotación errónea: se obtuvo %s, se esperaba 8.0 veces", indices[2].Valor)
		}
	})

	t.Run("Protección división por cero", func(t *testing.T) {
		// Si el pasivo es 0, no debe crashear la app, debe devolver 0.00
		indices := GenerarDashboard(100, 0, 0, 0, 0, 100, 0, 100, 10)

		if indices[0].Valor != "0.00" {
			t.Errorf("Error en división por cero: se obtuvo %s", indices[0].Valor)
		}
	})

	t.Run("Cálculo de Margen Neto", func(t *testing.T) {
		// Utilidad 25 sobre Ventas 100 = 25%
		indices := GenerarDashboard(0, 0, 0, 0, 0, 0, 0, 100, 25)

		// El margen neto es el índice 6 (índice 5 en el slice)
		margen := indices[5].Valor
		if margen != "25.00%" {
			t.Errorf("Margen Neto erróneo: se obtuvo %s, se esperaba 25.00%%", margen)
		}
	})
}
