// Package models Modelos para los Estados Financieros y la Base de datos
package models

type Cue struct {
	Codigo     int    `json:"codigo"`
	Nombre     string `json:"nombre"`
	Saldo      string `json:"saldo"`
	Categoria  string `json:"categoria"`
	EsCosto    bool   `json:"es_costo"`
	EsVariable bool   `json:"es_variable"`
	EsEfectivo bool   `json:"es_efectivo"`
}
type HtString struct {
	Nombre      string `json:"nombre"`
	Debe        string `json:"debe"`
	Haber       string `json:"haber"`
	Perdidas    string `json:"perdidas"`
	Ganancias   string `json:"ganancias"`
	Activo      string `json:"activo"`
	Pasivo      string `json:"pasivo"`
	ClassNombre string `json:"classnombre"`
	ClaTotal    string `json:"clatotal"`
	CostoDebe   string `json:"costodebe"`
	CostoHaber  string `json:"costohaber"`
}

type Ht struct {
	Nombre     string  `json:"nombre"`
	Debe       float64 `json:"debe"`
	Haber      float64 `json:"haber"`
	Perdidas   float64 `json:"perdidas"`
	Ganancias  float64 `json:"ganancias"`
	Activo     float64 `json:"activo"`
	Pasivo     float64 `json:"pasivo"`
	CostoDebe  float64 `json:"costodebe"`
	CostoHaber float64 `json:"costohaber"`
}
type Datos struct {
	Nombre string  `json:"nombre"`
	Saldo  string  `json:"saldo"`
	Valor  float64 `json:"valor"`
}

type Re struct {
	Nombre     string  `json:"nombre"`
	Col1       float64 `json:"col1"`
	Col2       float64 `json:"col2"`
	Col3       float64 `json:"col3"`
	ClasNombre string  `json:"clasnombre"`
	Cla1       string  `json:"cla1"`
	Cla2       string  `json:"cla2"`
	Cla3       string  `json:"cla3"`
}

type ReString struct {
	Nombre      string `json:"nombre"`
	Col1        string `json:"col1"`
	Col2        string `json:"col2"`
	Col3        string `json:"col3"`
	ClasNombre  string `json:"clasnombre"`
	Cla1        string `json:"cla1"`
	Cla2        string `json:"cla2"`
	Cla3        string `json:"cla3"`
	EsResultado bool   `json:"esresultado"`
}

type KR struct {
	Key   string
	Value ReString
}

type KV struct {
	Key   string
	Value HtString
}

type TotalesResultados struct {
	Ventas            float64
	CostoVentas       float64
	UtilidadBruta     float64
	GastosVenta       float64 // Los que tú llamas "de operación"
	GastosAdmin       float64
	TotalGastosOper   float64 // La suma de Venta + Admin
	UtilidadOperativa float64
	OtrosIngresos     float64
	OtrosGastos       float64
	UtilidadNeta      float64
	VentasNetas       float64
	MargenBruto       float64
	GastosFijos       float64
	GastosVariables   float64
}
