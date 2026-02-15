package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre   string `json:"nombre" form:"nombre"`
	Apellido string `json:"apellido" form:"apellido"`
	Email    string `json:"email" form:"email" gorm:"uniqueIndex"`
	Pase     string `json:"pase" form:"pase"`
	Role     string `json:"role" form:"role"`
	// --- NUEVOS CAMPOS ---
	SuscripcionActiva bool      `json:"suscripcion_activa" gorm:"default:false"`
	FechaFinPrueba    time.Time `json:"fecha_fin_prueba"` // Atajo para no calcular siempre
}

type Mensaje struct {
	gorm.Model
	UserID uint `json:"user_id"`
	//	User      User   `gorm:"foreignKey:UserID"`
	Consulta  string `json:"consulta" gorm:"type:text" form:"consulta"`
	Respuesta string `json:"respuesta" gorm:"type:text"`
	Estado    string `json:"estado" gorm:"default:'pendiente'"`
}

// En models/models.go

type CueDB struct {
	Codigo    int    `gorm:"primaryKey"`
	Nombre    string `gorm:"type:varchar(100)"`
	Saldo     string `gorm:"type:varchar(20)"`
	Categoria string `gorm:"type:varchar(50)"`
	EsCosto   bool   `gorm:"default:false"`
}

// TableName le dice a GORM que use el nombre exacto de la tabla que creamos

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
