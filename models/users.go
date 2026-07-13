package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre            string    `json:"nombre" form:"nombre"`
	Apellido          string    `json:"apellido" form:"apellido"`
	Email             string    `json:"email" form:"email" gorm:"uniqueIndex"`
	Pase              string    `json:"pase" form:"pase"`
	Role              string    `json:"role" form:"role"`
	SuscripcionActiva bool      `json:"suscripcion_activa" gorm:"default:false"`
	FechaFinPrueba    time.Time `json:"fecha_fin_prueba"`
	Verificado        bool      `json:"verificado" gorm:"default:false"`
	TokenVerificacion string    `json:"token_verificacion"`
}

type Mensaje struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User `gorm:"foreignKey:UserID" json:"-"`
	Consulta  string `json:"consulta" gorm:"type:text" form:"consulta"`
	Respuesta string `json:"respuesta" gorm:"type:text"`
	Estado    string `json:"estado" gorm:"default:'pendiente'"`
}

func (Mensaje) TableName() string {
	return "mensajes"
}

type CueDB struct {
	Codigo     int    `gorm:"primaryKey;column:codigo"`
	Nombre     string `gorm:"type:varchar(100);column:nombre"`
	Saldo      string `gorm:"type:varchar(20);column:saldo"`
	Categoria  string `gorm:"type:varchar(50);column:categoria"`
	EsCosto    bool   `gorm:"default:false;column:es_costo"`
	EsVariable bool   `gorm:"default:false;column:es_variable"` // <--- ESTE TAG ES VITAL
	EsEfectivo bool   `gorm:"default:false;column:es_efectivo"` // <--- ESTE TAMBIÉN
}

func (CueDB) TableName() string {
	return "nomenclatura"
}

type Transaccion struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	Monto      float64 `json:"monto"`
	Estado     string  `json:"estado"`     // "SUCCESS", "PENDING", "FAILED"
	Referencia string  `json:"referencia"` // El ID que te da QPayPro
	Pasarela   string  `json:"pasarela"`   // "qpaypro"
}

func (Transaccion) TableName() string {
	return "transacciones"
}
