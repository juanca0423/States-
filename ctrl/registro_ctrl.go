package ctrl

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"ef/db"
	"ef/models" // cambia por tu paquete
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// RegisterHandler crea un nuevo usuario con contraseña encriptada
func RegisterHandler(c *fiber.Ctx) error {
	type input struct {
		Nombre   string `json:"nombre"`
		Apellido string `json:"apellido"`
		Email    string `json:"email" validate:"required,email"`
		Pase     string `json:"pase" validate:"required,min=6"`
		Role     string `json:"role"` // opcional
	}

	var in input
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	// sanitización básica
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !emailRegex.MatchString(in.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email no válido",
		})
	}
	if len(in.Pase) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "La contraseña debe tener al menos 6 caracteres",
		})
	}

	// rol por defecto
	role := strings.ToLower(in.Role)
	switch role {
	case "admin", "cliente":
		// permitidos, pero podrías restringir "admin" a un endpoint separado
	default:
		role = "usuario"
	}

	// evitar duplicados
	var count int64
	db.DB.Model(&models.User{}).Where("email = ?", in.Email).Count(&count)
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "El email ya está registrado",
		})
	}

	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Pase), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al encriptar contraseña",
		})
	}

	// creación
	user := models.User{
		Nombre:   in.Nombre,
		Apellido: in.Apellido,
		Email:    in.Email,
		Pase:     string(hash),
		Role:     in.Role,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear el usuario",
		})
	}

	// respuesta sin el hash
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       user.ID,
		"nombre":   user.Nombre,
		"apellido": user.Apellido,
		"email":    user.Email,
		"role":     user.Role,
	})
}
