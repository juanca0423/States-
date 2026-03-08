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

	for _, c := range plantilla {
		valorReal := 0.0
		codigo := c.Codigo
		// Si es la cuenta de Ganancia, usamos el valor que calculamos en Resultados
		if codigo == 131002 {
			valorReal = utilidadNeta
		} else {
			// Si no, buscamos en el mapa normal
			if cuenta, existe := saldos[strconv.Itoa(c.Codigo)]; existe {
				valorReal = cuenta.Saldo
			}
		}

		// REGLA DE SIGNOS:
		// En el Pasivo/Patrimonio, las cuentas "Activo" RESTAN (como la Cuenta Personal)
		if c.Saldo == "activo" {
			suma -= valorReal
		} else {
			suma += valorReal
		}

		if valorReal != 0 {
			filas = append(filas, models.BaString{
				Codigo: strconv.Itoa(codigo),
				Nombre: c.Nombre,
				Col1:   config.FCont(valorReal),
			})
		}
	}

	*acumPat += suma
	filas = append(filas, models.BaString{
		Nombre: "Total " + titulo,
		Col3:   config.FCont(suma), // El total va a la columna 3
		Cla3:   "border-top",
	})
	return filas
}

func GenerarBalanceVista(plantilla []models.Cue, tituloGrupo string, acumuladorTotal *float64, saldosReales map[string]models.Cuenta, esPasivo bool, t *models.TotalesBalance) []models.BaString {
	var filas []models.BaString
	var sumaSubGrupo float64

	filas = append(filas, models.BaString{Nombre: tituloGrupo, ClasNombre: "ps-4 fw-bold text-muted small"})

	for _, c := range plantilla {
		valorReal := 0.0
		if cuenta, existe := saldosReales[strconv.Itoa(c.Codigo)]; existe {
			valorReal = cuenta.Saldo
		}

		// --- CORRECCIÓN DE LA LÓGICA DE INTERCEPCIÓN ---
		// --- CORRECCIÓN DE LA LÓGICA DE INTERCEPCIÓN ---
		if c.Codigo == 111301 {
			// 1. Guardamos el valor que venía de saldos reales como INICIAL
			t.InventarioInicial = valorReal

			// 2. Buscamos el inventario final para el Balance General
			if invFin, existe := saldosReales["220005"]; existe {
				valorReal = invFin.Saldo
				t.Inventario = valorReal // Este es el FINAL
			}
		}

		// LÓGICA DE SUMA O RESTA SEGÚN SECCIÓN (Mantenemos tu lógica original)
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
		}
	}

	*acumuladorTotal += sumaSubGrupo

	// Añadimos el subtotal del grupo
	filas = append(filas, models.BaString{
		Nombre: "Total " + tituloGrupo,
		Col2:   config.FCont(sumaSubGrupo),
		Cla2:   "border-top",
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

	// Corriente
	vista = append(vista, models.BaString{Nombre: "CORRIENTE", ClasNombre: "fw-bold ps-2"})
	vista = append(vista, GenerarBalanceVista(config.Disponible, "Disponibilidades", &t.ActivoCorriente, datos, false, &t)...)
	vista = append(vista, GenerarBalanceVista(config.Exigible, "Exigible", &t.ActivoCorriente, datos, false, &t)...)
	vista = append(vista, GenerarBalanceVista(config.Realisable, "Realizable", &t.ActivoCorriente, datos, false, &t)...)

	// Inyectamos Total Corriente en Columna 3
	vista = append(vista, models.BaString{Nombre: "Total Activo Corriente", Col3: config.FCont(t.ActivoCorriente), Cla3: "border-top"})

	// No Corriente
	vista = append(vista, models.BaString{Nombre: "NO CORRIENTE", ClasNombre: "fw-bold ps-2 mt-2"})
	vista = append(vista, GenerarBalanceVista(config.PropPlanEqui, "Propiedad, Planta y Equipo", &t.ActivoNoCorriente, datos, false, &t)...)
	vista = append(vista, GenerarBalanceVista(config.GtoIntan, "Gastos Intangibles", &t.ActivoNoCorriente, datos, false, &t)...)
	vista = append(vista, GenerarBalanceVista(config.GtoDiferidos, "Inversiones y Otros", &t.ActivoNoCorriente, datos, false, &t)...)

	// Inyectamos Total No Corriente en Columna 3
	vista = append(vista, models.BaString{Nombre: "Total Activo No Corriente", Col3: config.FCont(t.ActivoNoCorriente), Cla3: "border-top"})

	t.ActivoTotal = t.ActivoCorriente + t.ActivoNoCorriente
	vista = append(vista, models.BaString{Nombre: "SUMA DEL ACTIVO", Col4: config.FCont(t.ActivoTotal), ClasNombre: "fw-bold text-primary", Cla4: " border-bottom-double"})

	// --- 2. PASIVO ---
	vista = append(vista, models.BaString{Nombre: "PASIVO Y PATRIMONIO", ClasNombre: "fs-5 mt-4 text-danger"})

	// Pasivo Corriente
	vista = append(vista, GenerarBalanceVista(config.PasivoCorr, "Corto Plazo", &t.PasivoCorriente, datos, true, &t)...)
	vista = append(vista, models.BaString{Nombre: "Total Pasivo Corriente", Col3: config.FCont(t.PasivoCorriente), Cla3: "border-top"})

	// Pasivo No Corriente
	vista = append(vista, GenerarBalanceVista(config.PasivoNoCorr, "Largo Plazo", &t.PasivoNoCorriente, datos, true, &t)...)
	t.PasivoTotal = t.PasivoCorriente + t.PasivoNoCorriente
	vista = append(vista, models.BaString{Nombre: "Total Pasivo", Col3: config.FCont(t.PasivoTotal), Cla3: "border-top"})

	// --- 3. PATRIMONIO ---
	vista = append(vista, GenerarPatrimonioVista(config.PatriNeto, "Patrimonio Neto", &t.Patrimonio, datos, utilidadNeta)...)

	sumaPasivoPatrimonio := t.PasivoTotal + t.Patrimonio
	vista = append(vista, models.BaString{Nombre: "SUMA DEL PASIVO Y PATRIMONIO", Col4: config.FCont(sumaPasivoPatrimonio), ClasNombre: "text-danger", Cla4: "border-top"})

	PasPatrUtilidad := sumaPasivoPatrimonio + utilidadNeta
	vista = append(vista, models.BaString{Nombre: "GANANCIA ANTES DE IMPUESTOS Y RESERVAS", Col4: config.FCont(utilidadNeta), ClasNombre: "text-success-dark", Cla4: "border-bottom"})
	vista = append(vista, models.BaString{Nombre: "SUMA IGUAL AL PASIVO", Col4: config.FCont(PasPatrUtilidad), ClasNombre: "text-success-dark", Cla4: "border-bottom-double"})

	return vista, t
}
