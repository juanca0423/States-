// middleware/auth.go
package middleware

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims personalizados
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken crea un JWT firmado (úsalo en tu handler de login)
func GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// Require verifica que el usuario tenga *exactamente* el rol indicado
func Require(allowedRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := extractRole(c)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		if role != allowedRole {
			return fiber.NewError(fiber.StatusForbidden,
				fmt.Sprintf("se requiere rol %s", allowedRole))
		}
		return c.Next()
	}
}

// Any verifica que el usuario tenga *alguno* de los roles listados
func Any(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, err := extractRole(c)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		siExiste := slices.Contains(roles, userRole)
		if siExiste{
			return c.Next()
		}
		return fiber.NewError(fiber.StatusForbidden,
			"sin permisos suficientes")
	}
}

// extractRole valida el JWT y devuelve el rol
func extractRole(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	if auth == "" {
		return "", errors.New("falta header Authorization")
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("formato Authorization inválido")
	}
	tokenStr := parts[1]

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("método inesperado")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("token inválido")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("claims inválidos")
	}
	// Guardamos datos útiles para handlers posteriores
	c.Locals("uid", claims.UserID)
	c.Locals("role", claims.Role)
	return claims.Role, nil
}
