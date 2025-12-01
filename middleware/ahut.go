package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Check for authentication credentials
	// If credentials are valid, proceed to the next middleware or route
	// If credentials are invalid, return an unauthorized response
	return
}

func AdminMiddleware(c *fiber.Ctx) error {
	userRole := getRolUserFromContext(c)
	if userRole != AdminRole {
		return c.Status(fiber.StatusForbidden).SendString("Permission Denied")
  }
  return c.Next()
}

// GenerateJWT generates a JWT for the user
func GenerateJWT(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
    "email": user.Email,
    "role":  user.Role,
	})
	signedToken, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func AuthMiddle(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method and return the secret key
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
	}
	// Extract user information from the token and store it in the context
	c.Locals("user", getUserFromToken(token))
	return c.Next()
}
