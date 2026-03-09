package help

// ResumenCostos contiene los totales para el Dashboard y el Render
// En ef/help/costos_help.go o donde definiste el struct
type ResumenCostos struct {
	InvInicialMP    float64
	ComprasMP       float64
	GastosMP        float64
	DevMP           float64
	ComprasNetasMP  float64
	MPDisponible    float64
	InvFinalMP      float64
	MPConsumida     float64
	MOD             float64
	CostoPrimo      float64
	CIF             float64
	CostoProduccion float64
	InvInicialProc  float64
	InvFinalProc    float64
	CostoArtTerm    float64
	InvInicialPT    float64
	InvFinalPT      float64
	CostoVentas     float64
}

func CalcularCostosIndustriales(balanceRaw map[string]float64) ResumenCostos {
	getV := func(cod string) float64 { return balanceRaw[cod] }

	// --- Lógica extraída de tu controlador ---
	invInicialMP := getV("111303")
	comprasNetas := (getV("310200") + getV("310300")) + (getV("310400") + getV("310500") + getV("310600")) - getV("310700")
	mpUtilizada := (invInicialMP + comprasNetas) - getV("310800")

	mod := getV("320100") + getV("320200") + getV("320300")
	cif := getV("330100") + getV("330200") + getV("330300") + getV("330400") + getV("330500")

	costoPrimo := mpUtilizada + mod
	costoPeriodo := costoPrimo + cif // Tu CostoProduccion

	// Artículos Terminados
	costoArtTerm := (costoPeriodo + getV("111304")) - getV("340201")
	// Costo de Ventas
	costoVentas := (costoArtTerm + getV("111305")) - getV("340202")

	return ResumenCostos{
		MPConsumida:     mpUtilizada,
		MOD:             mod,
		CIF:             cif,
		CostoPrimo:      costoPrimo,
		CostoProduccion: costoPeriodo,
		CostoArtTerm:    costoArtTerm,
		CostoVentas:     costoVentas,
	}
}
