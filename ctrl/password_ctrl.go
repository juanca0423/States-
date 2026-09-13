package ctrl

import (
	"os"
	"time"

	"ef/db"
	"ef/models"
	"ef/servicios"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretPassword = []byte(os.Getenv("JWT_SECRET"))

func GetForgotPassword(c *fiber.Ctx) error {
	return c.Render("forgot-password", fiber.Map{
		"Title": "Recuperar Contraseña",
	})
}

func PostForgotPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Render("forgot-password", fiber.Map{
			"Title":   "Recuperar Contraseña",
			"Success": "Si el correo está registrado, recibirás un enlace de recuperación.",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecretPassword)
	if err != nil {
		return c.Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "Error al generar el token de recuperación.",
		})
	}

	err = servicios.EnviarCorreoRecuperacion(user.Email, user.Nombre, tokenString)
	if err != nil {
		return c.Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "No se pudo enviar el correo de recuperación.",
		})
	}

	return c.Render("forgot-password", fiber.Map{
		"Title":   "Recuperar Contraseña",
		"Success": "Hemos enviado un enlace a tu correo para restablecer tu contraseña.",
	})
}

func GetResetPassword(c *fiber.Ctx) error {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		return c.Redirect("/forgot-password")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretPassword, nil
	})

	if err != nil || !token.Valid {
		return c.Status(400).Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "El enlace de recuperación es inválido o ha expirado.",
		})
	}

	return c.Render("reset-password", fiber.Map{
		"Title": "Restablecer Contraseña",
		"Token": tokenStr,
	})
}

func PostResetPassword(c *fiber.Ctx) error {
	tokenStr := c.FormValue("token")
	newPassword := c.FormValue("pase")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretPassword, nil
	})

	if err != nil || !token.Valid {
		return c.Status(400).Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "El enlace de recuperación es inválido o ha expirado.",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(400).Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "Token inválido.",
		})
	}

	email := claims["email"].(string)

	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(404).Render("forgot-password", fiber.Map{
			"Title": "Recuperar Contraseña",
			"Error": "Usuario no encontrado.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).Render("reset-password", fiber.Map{
			"Title": "Restablecer Contraseña",
			"Token": tokenStr,
			"Error": "Error al procesar la contraseña.",
		})
	}

	user.Pase = string(hashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).Render("reset-password", fiber.Map{
			"Title": "Restablecer Contraseña",
			"Token": tokenStr,
			"Error": "Error al actualizar la contraseña en la base de datos.",
		})
	}

	return c.Redirect("/loguin")
}
