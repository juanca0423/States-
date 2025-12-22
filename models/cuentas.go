package models

type Cue struct {
	Codigo int     `json:"codigo"`
	Nombre string  `json:"nombre"`
	Saldo  string  `json:"Saldo"`
	Clases string  `json:"clases"`
	Valor  float64 `json:"valor"`
}

type HtString struct {
	Nombre    string `json:"nombre"`
	Debe      string `json:"debe"`
	Haber     string `json:"haber"`
	Perdidas  string `json:"perdidas"`
	Ganancias string `json:"ganancias"`
	Activo    string `json:"activo"`
	Pasivo    string `json:"pasivo"`
}

type Ht struct {
	Nombre    string  `json:"nombre"`
	Debe      float64 `json:"debe"`
	Haber     float64 `json:"haber"`
	Perdidas  float64 `json:"perdidas"`
	Ganancias float64 `json:"ganancias"`
	Activo    float64 `json:"activo"`
	Pasivo    float64 `json:"pasivo"`
}
type Datos struct {
	Nombre string  `json:"nombre"`
	Saldo  string  `json:"saldo"`
	Valor  float64 `json:"valor"`
}
