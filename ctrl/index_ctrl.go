package ctrl

import (
	"ef/db"
	"ef/middleware"
	"ef/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Index",
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
		"Title": "registro de usuarios",
	})
}

func PostRegistro(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "usuario registrado",
	})
}

func GetLayout(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func PostHandleLogin(c *fiber.Ctx) error {
	type input struct {
		Email string `json:"email"`
		Pase  string `json:"pase"`
	}
	var in input
	in.Email = c.FormValue("email")
	in.Pase = c.FormValue("pase")
	if in.Email == "" {
		return c.Render("index", fiber.Map{
			"Title": "el email esta vasio",
		})
	}
	/*
		if err := c.BodyParser(&in); err != nil {
			return fiber.ErrBadRequest
		}*/
	var u models.User
	if err := db.DB.Where("email = ?", in.Email).First(&u).Error; err != nil {
		return fiber.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pase), []byte(in.Pase)); err != nil {
		return fiber.ErrUnauthorized
	}
	tok, _ := middleware.GenerateToken(u.ID, u.Role)
	return c.JSON(fiber.Map{"token": tok})
}

func GetPerfil(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Que Hacemos",
	})
}
func GetPedCli(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Que Hacemos",
	})
}
func GetAdmDash(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Que Hacemos",
	})
}
