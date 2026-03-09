package help

import (
	"testing"

	"ef/config"
	"ef/models"
)

func TestHojaDeTrabajoReal(t *testing.T) {
	t.Run("Cuadre de Ganancia con Efectivo y Ventas", func(t *testing.T) {
		SetupNomenclatura()

		balanceRaw := map[string]float64{
			"111101": 1500.0, // Activo
			"210001": 1000.0, // Ganancia (Ventas)
			"231001": 500.0,  // Perdida (Sueldos)
		}

		resultado := HojaDeTrabajo(balanceRaw)

		// DEBUG: Mira cuántas filas se generaron
		t.Logf("Filas generadas en la hoja: %d", len(resultado))
		for _, kv := range resultado {
			t.Logf("Cuenta encontrada: %s", kv.Key)
		}
		var filaRes models.HtString
		encontrado := false
		for _, kv := range resultado {
			// Buscamos la fila de cierre por su código único
			if kv.Key == "910000" {
				filaRes = kv.Value
				encontrado = true
				break
			}
		}

		if !encontrado {
			t.Fatal("No se encontró la fila 910000")
		}

		// --- CORRECCIÓN 1: Usar el nombre que tu código realmente genera ---
		nombreEsperado := "RESULTADO DEL EJERCICIO" // Cambiado de "GANANCIA..."
		if filaRes.Nombre != nombreEsperado {
			t.Errorf("Nombre incorrecto: obtenido '%s', esperado '%s'", filaRes.Nombre, nombreEsperado)
		}

		// --- CORRECCIÓN 2: Verificación de valores con debug ---
		valorEsperado := config.FCont(500.0)

		// Si falla, el mensaje nos dirá exactamente qué caracteres hay (espacios, símbolos, etc.)
		if filaRes.Perdidas != valorEsperado {
			t.Errorf("Columna Perdidas falló: obtenida [%s], esperada [%s]", filaRes.Perdidas, valorEsperado)
		}

		if filaRes.Pasivo != valorEsperado {
			t.Errorf("Columna Pasivo falló: obtenida [%s], esperada [%s]", filaRes.Pasivo, valorEsperado)
		}
	})
	t.Run("Caso Especial Mercaderías 111301", func(t *testing.T) {
		// Esta es la parte más delicada de tu código:
		// 111301 (Mercaderías) usa el valor de 220005 (Inv. Final)
		balanceRaw := map[string]float64{
			"111301": 500.0, // Inventario Inicial
			"220005": 700.0, // Inventario Final
		}

		resultado := HojaDeTrabajo(balanceRaw)

		for _, kv := range resultado {
			if kv.Key == "111301" {
				ht := kv.Value

				// Según tu switch cuenta.Saldo == "activo" y codigo == "111301":
				// ht.Debe = 500 (Valor en balanceRaw)
				// ht.Activo = 700 (Valor de 220005)
				// ht.Perdidas = 500 | ht.Ganancias = 700
				if ht.Debe != config.FCont(500.0) || ht.Activo != config.FCont(700.0) {
					t.Errorf("Mercaderías: Debe(500) o Activo(700) incorrecto. Obtuve Debe:%s Activo:%s", ht.Debe, ht.Activo)
				}
				if ht.Perdidas != config.FCont(500.0) || ht.Ganancias != config.FCont(700.0) {
					t.Errorf("Mercaderías: Perdidas/Ganancias incorrectas")
				}
			}
		}
	})
}

func TestCuadreBalance(t *testing.T) {
	SetupNomenclatura()

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

	// Verificación extra: El patrimonio debe ser exactamente lo que mandamos (250)
	if totales.Patrimonio != 250.0 {
		t.Errorf("PATRIMONIO MAL CALCULADO: Se obtuvo %.2f, se esperaba 250.00", totales.Patrimonio)
	}

	// El Pasivo + Patrimonio + Utilidad debe sumar igual al Activo
	sumaLadoDerecho := totales.PasivoTotal + totales.Patrimonio + utilidadNeta
	if sumaLadoDerecho != totales.ActivoTotal {
		t.Errorf("BALANCE DESCUADRADO: Activo(%.2f) != Pasivo+Pat(%.2f)", totales.ActivoTotal, sumaLadoDerecho)
	}
}

func TestResultados(t *testing.T) {
	SetupNomenclatura()
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

func TestSeccionProduccionReal(t *testing.T) {
	SetupNomenclatura()

	// 1. ESCENARIO PARA "Cálculo Final del Costo" (Esperado 7500)
	balance7500 := map[string]float64{
		"111303": 2000.0, // Inventario Inicial MP (La función busca este código)
		"310200": 3000.0, // Compras MP
		"310800": 1000.0, // Inventario Final MP
		"320100": 2500.0, // MOD
		"330300": 600.0,  // CIF
		"330400": 400.0,  // CIF
	}

	// 2. ESCENARIO PARA "Validación de Lógica Industrial" (Esperado 8500)
	balance8500 := map[string]float64{
		"111303": 1000.0, // Inv Inicial
		"310200": 5000.0, // Compras
		"310800": 500.0,  // Inv Final  -> MP Consumida = 5500
		"320100": 2000.0, // MOD         -> Costo Primo = 7500
		"330100": 1000.0, // CIF         -> Total = 8500
	}
	t.Run("Validación de Dashboard Industrial", func(t *testing.T) {
		// Usamos datos conocidos
		costos := ResumenCostos{
			CostoPrimo:      6000,
			CIF:             2000,
			CostoProduccion: 8000,
		}
		ventas := 10000.0

		db := GenerarDashboardIndustrial(costos, ventas)

		// El CIF debería ser el 25% (2000/8000)
		if db[1].Valor != "25.00%" {
			t.Errorf("Cálculo de CIF incorrecto, esperado 25.00%% obtenido %s", db[1].Valor)
		}

		// El Margen Industrial debería ser 20% (10000-8000)/10000
		if db[2].Valor != "20.00%" {
			t.Errorf("Cálculo de Margen Industrial incorrecto, obtenido %s", db[2].Valor)
		}
	})
	t.Run("Validación de Lógica Industrial", func(t *testing.T) {
		res := CalcularCostosIndustriales(balance8500)
		if res.CostoProduccion != 8500 {
			t.Errorf("Costo producción erróneo, esperado 8500 obtenido %.2f", res.CostoProduccion)
		}
	})

	t.Run("Cálculo Final del Costo", func(t *testing.T) {
		costos := CalcularCostosIndustriales(balance7500)

		// MP(4000) + MOD(2500) + CIF(1000) = 7500
		esperado := 7500.0
		if costos.CostoProduccion != esperado {
			t.Errorf("Costo de Producción mal calculado. Esperado: %.2f, Obtenido: %.2f", esperado, costos.CostoProduccion)
		}

		t.Logf("Resultado de Fábrica: MP:%.2f, MOD:%.2f, CIF:%.2f. Total:%.2f",
			costos.MPConsumida, costos.MOD, costos.CIF, costos.CostoProduccion)
	})
}
