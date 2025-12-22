package ctrl

import (
	"ef/config"
	"ef/models"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type KV struct {
	Key   string
	Value models.HtString
}

func GenEstados(c *fiber.Ctx) error {
	// 1. fuerza el parse de multipart (Fiber lo hace solo si llamas a c.BodyParser,
	//    pero por claridad lo forzamos aquí)
	form, err := c.Request().MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "multipart error"})
	}
	// 2. form.Value es map[string][]string
	out := make(map[string]float64)
	for k, vv := range form.Value {
		if len(vv) > 0 {
			if f, err := strconv.ParseFloat(vv[0], 64); err == nil {
				out[k] = f
			}
		}
	}
	Balance := make(map[string]float64)
	for k, v := range out {
		if v != 0 {
			Balance[k] = v
		}
	}

	nomenclatura := config.CreaMap()
	// HoTra := make(map[string]models.Ht)
	var sumas models.Ht
	key := make([]KV, 0, len(Balance))
	sumas.Nombre = "Sumas"
	for k, v := range Balance {
		var dato models.Datos
		dato.Valor = v
		dato.Nombre = nomenclatura[k].Nombre
		dato.Saldo = nomenclatura[k].Saldo
		var ht models.HtString
		switch dato.Saldo {
		case "activo":
			if k == "111301" {
				ht.Nombre = dato.Nombre
				ht.Debe = config.FCont(dato.Valor)
				ht.Activo = config.FCont(Balance["220005"])
				ht.Perdidas = config.FCont(dato.Valor)
				ht.Ganancias = config.FCont(Balance["220005"])
				sumas.Debe += dato.Valor
				sumas.Activo += Balance["220005"]
				sumas.Perdidas += dato.Valor
				sumas.Ganancias += Balance["220005"]
				//HoTra[k] = ht
				key = append(key, KV{Key: k, Value: ht})
			} else {
				ht.Nombre = dato.Nombre
				ht.Debe = config.FCont(dato.Valor)
				ht.Activo = config.FCont(dato.Valor)
				sumas.Debe += dato.Valor
				sumas.Activo += dato.Valor
				//HoTra[k] = ht
				key = append(key, KV{Key: k, Value: ht})
			}
		case "pasivo":
			ht.Nombre = dato.Nombre
			ht.Haber = config.FCont(dato.Valor)
			ht.Pasivo = config.FCont(dato.Valor)
			sumas.Haber += dato.Valor
			sumas.Pasivo += dato.Valor
			//HoTra[k] = ht
			key = append(key, KV{Key: k, Value: ht})
		case "perdida":
			if k != "220001" {
				ht.Nombre = dato.Nombre
				ht.Debe = config.FCont(dato.Valor)
				ht.Perdidas = config.FCont(dato.Valor)
				sumas.Debe += dato.Valor
				sumas.Perdidas += dato.Valor
				//HoTra[k] = ht
				key = append(key, KV{Key: k, Value: ht})
			}
		case "ganancia":
			if k != "220005" {
				ht.Nombre = dato.Nombre
				ht.Haber = config.FCont(dato.Valor)
				ht.Ganancias = config.FCont(dato.Valor)
				sumas.Haber += dato.Valor
				sumas.Ganancias += dato.Valor
				//HoTra[k] = ht
				key = append(key, KV{Key: k, Value: ht})
			}
		default:
		}
	}
	//HoTra["900000"] = sumas
	ganancia := sumas.Ganancias - sumas.Perdidas
	var diferencia models.HtString
	var suma models.HtString
	suma.Nombre = sumas.Nombre
	suma.Debe = config.FCont(sumas.Debe)
	suma.Haber = config.FCont(sumas.Haber)
	suma.Perdidas = config.FCont(sumas.Perdidas)
	suma.Ganancias = config.FCont(sumas.Ganancias)
	suma.Activo = config.FCont(sumas.Activo)
	suma.Pasivo = config.FCont(sumas.Pasivo)

	key = append(key, KV{Key: "900000", Value: suma})
	if ganancia > 0 {
		diferencia.Nombre = "Ganancia"
		diferencia.Perdidas = config.FCont(ganancia)
		diferencia.Pasivo = config.FCont(ganancia)
		//HoTra["910000"] = diferencia
		key = append(key, KV{Key: "910000", Value: diferencia})
	} else {
		diferencia.Nombre = "Perdida"
		diferencia.Ganancias = config.FCont((-ganancia))
		diferencia.Activo = config.FCont((-ganancia))
		//HoTra["910000"] = diferencia
		key = append(key, KV{Key: "910000", Value: diferencia})
	}

	sort.Slice(key, func(i, j int) bool {
		return key[i].Key < key[j].Key
	})

	return c.Render("estados", fiber.Map{
		"keys": key,
		//"HojaTrabajo": HoTra,
	})
}

// fmt.Sprintf("%,.2f", monto)
