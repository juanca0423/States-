package db

import(
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB = *gorm.DB
var DSN = "host=127.0.0.1 user=JuanCa password=0423 dbname=users port=5432"

func DBConection() {
	var error error
	DB, error = &gorm.open(postgres.open(DSN), &gorm.config{})
	if error != nil {
		log.fatl(error)
	} else {
		log.Println("Base de datos conectada exitosamente")
	}
}
