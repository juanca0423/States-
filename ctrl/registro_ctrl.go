package ctrl

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"ef/db"
	"ef/models"
)

func VerificarCaptcha(response string) bool {
	secretKey := os.Getenv("EEFFS_APP") // Mejor si viene de os.Getenv("RECAPTCHA_SECRET")
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	resp, err := http.PostForm(apiURL, url.Values{
		"secret":   {secretKey},
		"response": {response},
	})
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Success
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func RegisterHandler(c *fiber.Ctx) error {
	captchaResponse := c.FormValue("g-recaptcha-response")
	if !VerificarCaptcha(captchaResponse) {
		return c.Render("register", fiber.Map{
			"Error": "Por favor, completa el captcha correctamente.",
		})
	}

	// 1. Definición del input
	type input struct {
		Nombre   string `json:"nombre" form:"nombre"`
		Apellido string `json:"apellido" form:"apellido"`
		Email    string `json:"email" form:"email"`
		Pase     string `json:"pase" form:"pase"`
		Role     string `json:"role" form:"role"`
	}

	var in input
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "JSON o Formulario inválido",
		})
	}

	// 2. Sanitización y validación
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !emailRegex.MatchString(in.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Email no válido",
		})
	}
	if len(in.Pase) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "La contraseña debe tener al menos 6 caracteres",
		})
	}

	// 3. Lógica de Roles
	role := strings.ToLower(in.Role)
	if role != "admin" && role != "cliente" {
		role = "usuario"
	}

	// 4. Evitar duplicados
	var count int64
	db.DB.Model(&models.User{}).Where("email = ?", in.Email).Count(&count)
	if count > 0 {
		return c.Render("register", fiber.Map{
			"Error": "Este correo ya está en uso, intenta con otro.",
			"Datos": in,
		})
	}

	// 5. Encriptar contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Pase), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error al encriptar contraseña",
		})
	}

	// 6. Creación del usuario en la BD
	user := models.User{
		Nombre:            in.Nombre,
		Apellido:          in.Apellido,
		Email:             in.Email,
		Pase:              string(hash),
		Role:              role, // Usamos la variable 'role' ya validada
		SuscripcionActiva: false,
		FechaFinPrueba:    time.Now().AddDate(0, 0, 30),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Render("register", fiber.Map{
			"Error": "Error al guardar en BD",
		})
	}

	// 7. Respuesta
	return c.Render("bienvenida", fiber.Map{
		"Nombre": user.Nombre,
		"Email":  user.Email,
	})
}
