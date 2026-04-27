package help

/*
func Resultados1(Balance map[string]float64) []models.KR {
	Res := make([]models.KR, 0, len(Balance))
	ing, ving := RecCol3(Balance, config.Ingresos)
	Res = append(Res, ing...)
	if ving > 0 {
		var dato models.ReString
		dato.Col3 = config.FCont(ving)
		dato.Nombre = msg.Msg["es"]["Vn"]
		Res = append(Res, models.KR{Key: "210010", Value: dato})
	}
	invini, vinvini := InvIni(Balance)
	Res = append(Res, invini...)
	comp, vcomp := Comp(Balance)
	Res = append(Res, comp...)
	cosven := CostoVentas(Balance, vinvini, vcomp)
	Res = append(Res, cosven...)

	gtoven, vgtoven := RecCol1(Balance, config.GtoVentas)
	Res = append(Res, gtoven...)

	gtoadm, acugto := RecCol1Tot(Balance, config.GtoAdmin, vgtoven)
	Res = append(Res, gtoadm...)
	var dato models.ReString
	dato.Nombre = msg.Msg["es"]["Mb"]
	dato.Col3 = config.FCont(ving - acugto)
	Res = append(Res, models.KR{Key: "232900", Value: dato})

	return Res
}

func Resultadospropuesta(Balance map[string]float64) []models.KR {
	Res := make([]models.KR, 0)

	// --- SECCIÓN 1: INGRESOS ---
	filasIng, vVentasNetas := RecCol3(Balance, config.Ingresos)
	Res = append(Res, filasIng...)

	// --- SECCIÓN 2: COSTO DE VENTAS ---
	filasInvIni, vInvIni := InvIni(Balance)
	Res = append(Res, filasInvIni...)

	filasComp, vCompNetas := Comp(Balance)
	Res = append(Res, filasComp...)

	filasCV, vCostoVentas := CostoVentas(Balance, vInvIni, vCompNetas)
	Res = append(Res, filasCV...)

	// --- MARGEN BRUTO ---
	vUtilidadBruta := vVentasNetas - vCostoVentas
	Res = append(Res, models.KR{Key: "220400", Value: models.ReString{
		Nombre: "MARGEN BRUTO (UTILIDAD BRUTA)",
		Col3:   config.FCont(vUtilidadBruta),
		Cla3:   "total",
	}})

	// --- SECCIÓN 3: GASTOS DE OPERACIÓN ---
	// Usamos Col1 y Col2 para los gastos
	gtoven, vSumGtoVen := RecCol1(Balance, config.GtoVentas)
	Res = append(Res, gtoven...)

	gtoadm, vSumGtoAdm := RecCol1(Balance, config.GtoAdmin)
	Res = append(Res, gtoadm...)

	vTotalGastos := vSumGtoVen + vSumGtoAdm
	Res = append(Res, models.KR{Key: "233000", Value: models.ReString{
		Nombre: "TOTAL GASTOS DE OPERACIÓN",
		Col2:   config.FCont(vTotalGastos),
		Col3:   config.FCont(vTotalGastos), // Lo pasamos a col3 para restar
		Cla3:   "total",
	}})

	// --- RESULTADO DE OPERACIÓN ---
	vResultadoOp := vUtilidadBruta - vTotalGastos
	Res = append(Res, models.KR{Key: "240000", Value: models.ReString{
		Nombre: "RESULTADO DE OPERACIÓN",
		Col3:   config.FCont(vResultadoOp),
		Cla3:   "doble-linea", // Para resaltar el final
	}})

	return Res
}*/
