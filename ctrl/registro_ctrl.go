package ctrl

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	"ef/servicios" // Asegúrate de que el path a tu carpeta servicios sea correcto
)

// Función auxiliar para generar un token seguro
func generarToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

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

func VerificarCuenta(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).SendString("Token inválido")
	}

	var usuario models.User
	// Buscamos al usuario por el token
	result := db.DB.Where("token_verificacion = ?", token).First(&usuario)

	if result.Error != nil {
		return c.Status(404).Render("error", fiber.Map{
			"Mensaje": "El enlace de verificación ha expirado o es incorrecto.",
		})
	}

	// Actualizamos al usuario
	usuario.Verificado = true
	usuario.TokenVerificacion = "" // Limpiamos el token por seguridad
	db.DB.Save(&usuario)

	return c.Render("bienvenida", fiber.Map{
		"Nombre":  usuario.Nombre,
		"Email":   usuario.Email,
		"Mensaje": "Tu cuenta ha sido verificada. Ya puedes disfrutar de EEFFS.",
	})
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func RegisterHandler(c *fiber.Ctx) error {
	captchaResponse := c.FormValue("g-recaptcha-response")
	if !VerificarCaptcha(captchaResponse) {
		return c.Render("register", fiber.Map{
			"Error": "Por favor, completa el captcha correctamente.",
		})
	}

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

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !emailRegex.MatchString(in.Email) {
		return c.Render("register", fiber.Map{"Error": "Email no válido", "Datos": in})
	}
	if len(in.Pase) < 6 {
		return c.Render("register", fiber.Map{"Error": "La contraseña debe tener al menos 6 caracteres", "Datos": in})
	}

	var count int64
	db.DB.Model(&models.User{}).Where("email = ?", in.Email).Count(&count)
	if count > 0 {
		return c.Render("register", fiber.Map{
			"Error": "Este correo ya está en uso, intenta con otro.",
			"Datos": in,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Pase), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Error al encriptar"})
	}

	// 1. Generamos el token único para este usuario
	token := generarToken()

	// 2. Creamos el usuario con Verificado: false
	user := models.User{
		Nombre:            in.Nombre,
		Apellido:          in.Apellido,
		Email:             in.Email,
		Pase:              string(hash),
		Role:              "usuario", // Forzamos usuario por seguridad en registro público
		SuscripcionActiva: false,
		FechaFinPrueba:    time.Now().AddDate(0, 0, 30),
		Verificado:        false, // Campo nuevo
		TokenVerificacion: token, // Campo nuevo
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Render("register", fiber.Map{"Error": "Error al guardar en BD"})
	}

	// 3. ENVIAR CORREO EN SEGUNDO PLANO
	// Usamos 'go' para que el usuario no tenga que esperar a que Resend responda
	go func() {
		err := servicios.EnviarCorreoVerificacion(user.Email, user.Nombre, token)
		if err != nil {
			fmt.Printf("❌ Error enviando correo a %s: %v\n", user.Email, err)
		}
	}()

	// 4. Cambiamos la respuesta. Ya no va a 'bienvenida' directo,
	// sino a una página que le pide revisar su correo.
	return c.Render("confirmacion_enviada", fiber.Map{
		"Email":  user.Email,
		"Nombre": user.Nombre,
	})
}
