package ctrl

import (
	"strings"
	"time"

	"ef/db"
	"ef/middleware"
	"ef/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "EE FF",
	})
}

func GetLoguin(c *fiber.Ctx) error {
	return c.Render("loguin", fiber.Map{
		"Title": "Loguin",
	})
}

func GetAbout(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Que Hacemos",
	})
}

func GetRejistro(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Registrate es Gratis",
	})
}

func GetPlanes(c *fiber.Ctx) error {
	reason := c.Query("reason")
	mensaje := "Para continuar usando States debes suscribirte o renovar tu plan."
	if reason == "expired" {
		mensaje = "Tu periodo de prueba ha expirado. Elige un plan para continuar usando la aplicación."
	}

	return c.Render("planes", fiber.Map{
		"Title":       "Suscripción requerida",
		"Mensaje":     mensaje,
		"EsSuscrito":  c.Locals("Suscrito"),
		"TextoDias":   c.Locals("TextoDias"),
		"ClaseAlerta": c.Locals("ClaseAlerta"),
	})
}

func GetLayout(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func PostHandleLogin(c *fiber.Ctx) error {
	type input struct {
		Email string `form:"email"`
		Pase  string `form:"pase"`
	}

	var in input
	if err := c.BodyParser(&in); err != nil {
		return c.Render("loguin", fiber.Map{"Error": "Datos inválidos"})
	}

	// 1. Buscar usuario
	var u models.User
	emailLimpio := strings.ToLower(strings.TrimSpace(in.Email))

	if err := db.DB.Where("email = ?", emailLimpio).First(&u).Error; err != nil {
		return c.Render("loguin", fiber.Map{
			"Error": "Correo o contraseña no válidos",
			"Email": emailLimpio,
		})
	}

	// 2. Validar Contraseña (IMPORTANTE: Validar contraseña ANTES que la verificación)
	// Así no le decimos a un hacker si la cuenta está verificada o no si no sabe la clave.
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pase), []byte(in.Pase)); err != nil {
		return c.Render("loguin", fiber.Map{
			"Error": "Correo o contraseña no válidos",
			"Email": emailLimpio,
		})
	}

	// 3. NUEVO: Verificar si la cuenta está activa (Usando 'u' no 'usuario')
	if !u.Verificado {
		return c.Render("loguin", fiber.Map{ // Asegúrate que el nombre sea "loguin" o "login" según tus archivos
			"Error": "Tu cuenta aún no ha sido verificada. Revisa tu correo electrónico.",
			"Email": emailLimpio,
		})
	}

	// 4. Generar Token
	tok, err := middleware.GenerateToken(u.ID, u.Role)
	if err != nil {
		return c.Render("error-sistema", fiber.Map{
			"Codigo":  "500",
			"Mensaje": "No pudimos generar tu sesión de auditoría.",
		})
	}

	// 5. Guardar en Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tok,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		Secure:   true, // Ya que estás en eeffs.com (HTTPS), cámbialo a true
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Redirect("/eeff")
}
