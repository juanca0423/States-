package help

import (
	"math"
	"strconv"

	"ef/config"
	"ef/models"
)

func GenerarPatrimonioVista(plantilla []models.Cue, titulo string, acumPat *float64, saldos map[string]models.Cuenta, utilidadNeta float64) []models.BaString {
	var filas []models.BaString
	var suma float64
	ultimoIdx := -1

	// 1. Procesamos las cuentas reales que SÍ existen en la nomenclatura (Capital, Reservas, etc.)
	for _, c := range plantilla {
		valorReal := 0.0
		if cuenta, existe := saldos[strconv.Itoa(c.Codigo)]; existe {
			valorReal = cuenta.Saldo
		}

		if c.Saldo == "activo" {
			suma -= valorReal // Aquí restarán las acciones en tesorería automáticamente
		} else {
			suma += valorReal
		}

		if valorReal != 0 {
			filas = append(filas, models.BaString{
				Codigo: strconv.Itoa(c.Codigo),
				Nombre: c.Nombre,
				Col1:   config.FCont(valorReal),
			})
			ultimoIdx = len(filas) - 1
		}
	}

	// 2. MOSTRAMOS la utilidad calculada por el sistema como fila informativa,
	// pero NO la sumamos al patrimonio aquí. El objetivo es que `Patrimonio`
	// represente solo las cuentas de capital; el resultado del ejercicio se
	// maneja por separado al cuadrar el balance.
	if utilidadNeta != 0 {
		filas = append(filas, models.BaString{
			Codigo: "SISTEMA", // Identificador claro de que es un cálculo automático
			Nombre: "Utilidad (o Pérdida) del Ejercicio",
			Col1:   config.FCont(utilidadNeta),
		})
		ultimoIdx = len(filas) - 1
	}

	// Cerramos la Col1 con su línea de subtotal
	if ultimoIdx != -1 {
		filas[ultimoIdx].Cla1 = "border-bottom"
	}

	// 3. Acumulamos al gran total del Pasivo + Patrimonio
	*acumPat += suma

	filas = append(filas, models.BaString{
		Nombre: "Total " + titulo,
		Col3:   config.FCont(suma),
		Cla3:   "border-top",
	})
	return filas
}

func GenerarBalanceVista(plantilla []models.Cue, tituloGrupo string, acumuladorTotal *float64, saldosReales map[string]models.Cuenta, esPasivo bool, t *models.TotalesBalance) []models.BaString {
	var filas []models.BaString
	var sumaSubGrupo float64

	filas = append(filas, models.BaString{Nombre: tituloGrupo, ClasNombre: "ps-4 fw-bold text-muted small"})

	// Variable para rastrear el índice de la última fila numérica agregada
	ultimoIndiceNumerico := -1

	for _, c := range plantilla {
		valorReal := 0.0
		if cuenta, existe := saldosReales[strconv.Itoa(c.Codigo)]; existe {
			valorReal = cuenta.Saldo
		}

		if c.Codigo == 111301 {
			t.InventarioInicial = valorReal
			if invFin, existe := saldosReales["220005"]; existe {
				valorReal = invFin.Saldo
				t.Inventario = valorReal
			}
		}

		if esPasivo {
			if c.Saldo == "activo" {
				sumaSubGrupo -= valorReal
			} else {
				sumaSubGrupo += valorReal
			}
		} else {
			if c.Saldo == "pasivo" {
				sumaSubGrupo -= valorReal
			} else {
				sumaSubGrupo += valorReal
			}
		}

		if valorReal != 0 {
			filas = append(filas, models.BaString{
				Codigo: strconv.Itoa(c.Codigo),
				Nombre: c.Nombre,
				Col1:   config.FCont(valorReal),
			})
			ultimoIndiceNumerico = len(filas) - 1 // Guardamos la posición
		}
	}

	// --- LA MAGIA: Aplicar línea a la última cuenta del parcial ---
	if ultimoIndiceNumerico != -1 {
		filas[ultimoIndiceNumerico].Cla1 = "border-bottom"
	}

	*acumuladorTotal += sumaSubGrupo

	filas = append(filas, models.BaString{
		Nombre: "Total " + tituloGrupo,
		Col2:   config.FCont(sumaSubGrupo),
		Cla2:   "border-top", // Esta es la línea que recibe el subtotal en Col2
	})

	return filas
}

func GenerarFilaCierre(totalActivo float64, totalPasivoPatrimonio float64) []models.BaString {
	var filas []models.BaString

	// Calculamos la diferencia técnica
	diferencia := math.Abs(totalActivo - totalPasivoPatrimonio)

	filas = append(filas, models.BaString{
		Nombre:     "TOTAL DEL ACTIVO",
		Col4:       config.FCont(totalActivo),
		ClasNombre: "fw-bold text-info",     // Un color para resaltar
		Cla4:       "border-double fw-bold", // <--- ACTIVAR DOBLE LÍNEA
	})

	// Fila de Pasivo + Patrimonio
	claseCuadre := "border-double fw-bold"
	mensaje := "TOTAL PASIVO Y PATRIMONIO"

	if diferencia > 0.01 {
		mensaje = "⚠️ BALANCE DESCUADRADO"
		claseCuadre = "border-double text-danger"
	} else {
		// Agregamos una fila extra de "Cuadre Verificado" (Opcional)
		filas = append(filas, models.BaString{
			Nombre:     "SITUACIÓN FINANCIERA:",
			Col1:       "CUADRADO ✔️",
			ClasNombre: "small text-muted italic",
			Cla1:       "text-success-dark",
		})
	}

	filas = append(filas, models.BaString{
		Nombre:     mensaje,
		Col4:       config.FCont(totalPasivoPatrimonio),
		ClasNombre: "fw-bold",
		Cla4:       claseCuadre,
	})

	return filas
}

func GenerarTodoElBalance(datos map[string]models.Cuenta, utilidadNeta float64) ([]models.BaString, models.TotalesBalance) {
	var vista []models.BaString
	var t models.TotalesBalance

	// --- 1. ACTIVO ---
	vista = append(vista, models.BaString{Nombre: "ACTIVO", ClasNombre: "fw-bold fs-5 text-success"})

	// Corriente (Solo se agregan si tienen datos)
	vista = append(vista, GenerarBalanceVista(config.Disponible, "Disponibilidades", &t.ActivoCorriente, datos, false, &t)...)

	vista = append(vista, GenerarBalanceVista(config.Exigible, "Exigible", &t.ActivoCorriente, datos, false, &t)...)

	vista = append(vista, GenerarBalanceVista(config.Realisable, "Realizable", &t.ActivoCorriente, datos, false, &t)...)

	if t.ActivoCorriente > 0 {
		vista = append(vista, models.BaString{Nombre: "Total Activo Corriente", Col3: config.FCont(t.ActivoCorriente), Cla3: "border-top"})
	}

	// No Corriente
	vista = append(vista, GenerarBalanceVista(config.ActivoNoCorr, "Activos No Corrientes", &t.ActivoNoCorriente, datos, false, &t)...)

	vista = append(vista, GenerarBalanceVista(config.PropPlanEqui, "Propiedad, Planta y Equipo", &t.ActivoNoCorriente, datos, false, &t)...)

	vista = append(vista, GenerarBalanceVista(config.GtoIntan, "Gastos Intangibles", &t.ActivoNoCorriente, datos, false, &t)...)

	vista = append(vista, GenerarBalanceVista(config.GtoDiferidos, "Gastos Diferidos", &t.ActivoNoCorriente, datos, false, &t)...)

	if t.ActivoNoCorriente > 0 {
		vista = append(vista, models.BaString{Nombre: "Total Activo No Corriente", Col3: config.FCont(t.ActivoNoCorriente), Cla3: "border-top"})
	}

	t.ActivoTotal = t.ActivoCorriente + t.ActivoNoCorriente
	// En la sección de Sumas Finales:
	vista = append(vista, models.BaString{
		Nombre:      "SUMA DEL ACTIVO",
		Col4:        config.FCont(t.ActivoTotal),
		ClasNombre:  "fw-bold text-primary",
		Cla4:        "border-double", // Cambiado de border-bottom-double
		EsResultado: true,
	})

	// --- 2. PASIVO ---
	vista = append(vista, models.BaString{Nombre: "PASIVO Y PATRIMONIO", ClasNombre: "fs-5 mt-4 text-danger"})

	// Pasivo Corriente
	vista = append(vista, GenerarBalanceVista(config.PasivoCorr, "Corto Plazo", &t.PasivoCorriente, datos, true, &t)...)

	if t.PasivoCorriente > 0 {
		vista = append(vista, models.BaString{Nombre: "Total Pasivo Corriente", Col3: config.FCont(t.PasivoCorriente), Cla3: "border-top"})
	}

	// Pasivo No Corriente
	vista = append(vista, GenerarBalanceVista(config.PasivoNoCorr, "Largo Plazo", &t.PasivoNoCorriente, datos, true, &t)...)
	t.PasivoTotal = t.PasivoCorriente + t.PasivoNoCorriente
	vista = append(vista, models.BaString{Nombre: "Total Pasivo", Col3: config.FCont(t.PasivoTotal), Cla3: "border-top"})

	// --- 3. PATRIMONIO ---
	// t.Patrimonio aquí ya se calcula completo, absorbiendo internamente la utilidadNeta
	vista = append(vista, GenerarPatrimonioVista(config.PatriNeto, "Patrimonio Neto", &t.Patrimonio, datos, utilidadNeta)...)

	// 1. Calculamos la suma real final
	t.PasivoPatrimonioTotal = t.PasivoTotal + t.Patrimonio

	// 2. Inyectamos tus filas de cierre (las que muestran si está "CUADRADO ✔️" o "DESCUADRADO")
	filasCierre := GenerarFilaCierre(t.ActivoTotal, t.PasivoPatrimonioTotal)
	vista = append(vista, filasCierre...)

	return vista, t
}
