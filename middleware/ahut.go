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

	// BLINDAJE DE SEGURIDAD: Evita que pasen estructuras con datos vacíos o en cero
	if claims.UserID == 0 || claims.Role == "" {
		return nil, errors.New("claims incompletos o nulos")
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
	c.Locals("UserRole", user.Role)
	c.Locals("uid", claims.UserID)
	// En AuthStatus guarda esto:
	c.Locals("Suscrito", user.SuscripcionActiva)
	c.Locals("FinPrueba", user.FechaFinPrueba)
	// Dentro de AuthStatus, después de obtener el usuario de la DB:
	dias := max(0, int(time.Until(user.FechaFinPrueba).Hours()/24))
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
	roleRaw := c.Locals("UserRole")
	if roleRaw == nil {
		roleRaw = c.Locals("role")
	}
	role := fmt.Sprint(roleRaw)

	uidRaw := c.Locals("uid")
	uid, ok := uidRaw.(uint)
	if !ok {
		return c.Redirect("/loguin")
	}

	if role == "admin" {
		return c.Next()
	}

	// Recuperamos con validación de existencia
	suscritoRaw := c.Locals("Suscrito")
	finPruebaRaw := c.Locals("FinPrueba")

	if suscritoRaw == nil || finPruebaRaw == nil {
		var user models.User
		if err := db.DB.Select("suscripcion_activa", "fecha_fin_prueba", "role").First(&user, uid).Error; err != nil {
			fmt.Printf("⚠️ Error al validar suscripción: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).Render("error-sistema", fiber.Map{
				"Codigo":  "500",
				"Mensaje": "No se pudo verificar tu suscripción. Inténtalo nuevamente en unos minutos.",
			})
		}
		suscritoRaw = user.SuscripcionActiva
		finPruebaRaw = user.FechaFinPrueba
		role = user.Role
		c.Locals("Suscrito", suscritoRaw)
		c.Locals("FinPrueba", finPruebaRaw)
		c.Locals("role", role)
		c.Locals("UserRole", role)
	}

	suscrito, _ := suscritoRaw.(bool)
	finPrueba, _ := finPruebaRaw.(time.Time)

	if suscrito {
		return c.Next()
	}

	// Acceso mientras la prueba no haya vencido (antes o el mismo día de fin).
	if !time.Now().After(finPrueba) {
		return c.Next()
	}

	return c.Redirect("/planes?reason=expired")
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

// SetJWTSecret permite cambiar la clave JWT en tiempo de ejecución (útil para tests)
func SetJWTSecret(secret string) {
	jwtKey = []byte(secret)
}
