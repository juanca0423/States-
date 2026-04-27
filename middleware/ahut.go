// Package middleware
package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"ef/db"
	"ef/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Usamos un solo nombre para la clave y los claims
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// --- UTILIDADES ---

func GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("error al obtener claims")
	}
	return claims, nil
}

// --- MIDDLEWARES DE CONTROL ---

// AuthStatus Se usa globalmente para el Menú (HBS)
func AuthStatus(c *fiber.Ctx) error {
	path := c.Path()
	// 1. SKIP EFICIENTE: Si es estático, no toques la DB
	if len(path) > 4 {
		ext := path[len(path)-4:]
		if ext == ".css" || ext == ".js" || ext == ".png" || ext == ".jpg" || ext == ".ico" {
			return c.Next()
		}
	}
	cookie := c.Cookies("jwt")
	if cookie == "" {
		c.Locals("IsLogged", false)
		return c.Next()
	}
	claims, err := ValidateToken(cookie)
	if err != nil {
		c.ClearCookie("jwt")
		c.Locals("IsLogged", false)
		return c.Next()
	}
	// 2. EVITAR CONSULTAS REPETIDAS:
	// Si ya tenemos los datos en Locals, no volvemos a consultar la DB
	if c.Locals("IsLogged") == true && c.Locals("UserName") != nil {
		return c.Next()
	}
	// 3. CONSULTA SEGURA:
	// Si la conexión global de la DB está caída por un reinicio de Air,
	// evitamos que el middleware rompa la petición.
	if db.DB == nil {
		fmt.Println("⚠️ Advertencia: Intento de AuthStatus sin conexión a DB")
		return c.Next()
	}
	var user models.User
	// Select optimizado para PgBouncer
	result := db.DB.Select("nombre", "role", "suscripcion_activa", "fecha_fin_prueba").
		First(&user, claims.UserID)
	if result.Error != nil {
		// Si el error es de conexión (DNS), no cerramos la sesión,
		// solo dejamos que pase sin datos de usuario para que no "explote" la web
		fmt.Printf("⚠️ Error de DB en AuthStatus: %v\n", result.Error)
		return c.Next()
	}
	c.Locals("IsLogged", true)
	c.Locals("UserName", user.Nombre)
	c.Locals("role", user.Role)
	c.Locals("uid", claims.UserID)
	// En AuthStatus guarda esto:
	c.Locals("Suscrito", user.SuscripcionActiva)
	c.Locals("FinPrueba", user.FechaFinPrueba)
	// Dentro de AuthStatus, después de obtener el usuario de la DB:
	dias := int(time.Until(user.FechaFinPrueba).Hours() / 24)
	c.Locals("EsTrial", !user.SuscripcionActiva)
	c.Locals("TextoDias", fmt.Sprintf("%d días restantes", dias))

	if dias < 5 {
		c.Locals("ClaseAlerta", "bg-danger")
	} else {
		c.Locals("ClaseAlerta", "bg-warning text-dark")
	}
	return c.Next()
}

// AuthRequired: Bloquea acceso si no hay sesión

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		if c.Method() == "GET" {
			return c.Redirect("/loguin")
		}
		return fiber.ErrUnauthorized
	}

	claims, err := ValidateToken(cookie)
	if err != nil {
		if c.Method() == "GET" {
			return c.Redirect("/loguin")
		}
		return fiber.ErrUnauthorized
	}

	c.Locals("uid", claims.UserID)
	c.Locals("role", claims.Role)
	return c.Next()
}

// CheckSubscription Verifica pago o trial
func CheckSubscription(c *fiber.Ctx) error {
	role := c.Locals("UserRole") // Usamos el nombre que seteamos en AuthStatus
	uid := c.Locals("uid")

	if role == "admin" {
		return c.Next()
	}
	if uid == nil {
		return c.Redirect("/loguin")
	}

	// Recuperamos con validación de existencia
	suscritoRaw := c.Locals("Suscrito")
	finPruebaRaw := c.Locals("FinPrueba")

	// Si por alguna razón los datos no están en Locals (ej. error de DB previo)
	// volvemos a verificar para no dejar pasar a alguien sin permiso o romper la app
	if suscritoRaw == nil || finPruebaRaw == nil {
		fmt.Println("⚠️ Datos de suscripción no encontrados en Locals, redirigiendo...")
		return c.Redirect("/loguin")
	}

	suscrito := suscritoRaw.(bool)
	finPrueba := finPruebaRaw.(time.Time)

	if suscrito || time.Now().Before(finPrueba) {
		return c.Next()
	}

	return c.Render("pago_requerido", fiber.Map{
		"Mensaje": "Tu periodo de prueba ha finalizado. Suscríbete para continuar.",
	})
}

// Require: Roles específicos (ej. Admin)

func Require(allowedRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != allowedRole {
			return fiber.NewError(fiber.StatusForbidden, "Acceso denegado")
		}
		return c.Next()
	}
}
