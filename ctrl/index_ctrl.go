package ctrl

import (
	"strings"
	"time"

	"ef/db"
	"ef/middleware"
	"ef/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	// Cambia la línea de la consulta por esta:
	if err := db.DB.Session(&gorm.Session{PrepareStmt: false}).Where("email = ?", emailLimpio).First(&u).Error; err != nil {
		// No revelamos si el correo existe o no, pero devolvemos el email para su comodidad
		return c.Render("loguin", fiber.Map{
			"Error": "Correo o contraseña no válidos",
			"Email": emailLimpio,
		})
	}

	// 2. Validar Contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pase), []byte(in.Pase)); err != nil {
		return c.Render("loguin", fiber.Map{
			"Error": "Correo o contraseña no válidos",
			"Email": emailLimpio,
		})
	}

	// 3. Generar Token
	tok, err := middleware.GenerateToken(u.ID, u.Role)
	if err != nil {
		return c.Render("error-sistema", fiber.Map{
			"Codigo":  "500",
			"Mensaje": "No pudimos generar tu sesión de auditoría.",
		})
	}

	// 4. Guardar en Cookie con configuración profesional
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tok,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // ¡Recuerda ponerlo en TRUE cuando subas a internet con HTTPS!
		SameSite: "Lax",
		Path:     "/", // Asegura que la cookie sea válida en todas las rutas del sistema
	})

	// 5. Redirigir al panel principal
	return c.Redirect("/eeff")
}
