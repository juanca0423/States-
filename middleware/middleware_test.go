package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ef/middleware"

	"github.com/gofiber/fiber/v2"
)

func TestGenerateAndValidateToken(t *testing.T) {
	middleware.SetJWTSecret("test-secret-123")

	tok, err := middleware.GenerateToken(42, "user")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	claims, err := middleware.ValidateToken(tok)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected uid 42, got %v", claims.UserID)
	}
	if claims.Role != "user" {
		t.Fatalf("expected role 'user', got %v", claims.Role)
	}
}

func TestRequireMiddleware(t *testing.T) {
	app := fiber.New()

	// Caso forbidden
	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("role", "user")
		return c.Next()
	}, middleware.Require("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}

	// Caso permitido
	app2 := fiber.New()
	app2.Get("/", func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return c.Next()
	}, middleware.Require("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	resp2, err := app2.Test(httptest.NewRequest("GET", "/", nil), 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp2.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp2.StatusCode)
	}
}

func TestAuthRequiredRedirectsWhenNoCookie(t *testing.T) {
	app := fiber.New()
	app.Get("/prot", middleware.AuthRequired, func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/prot", nil)
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}

	if resp.StatusCode != fiber.StatusFound && resp.StatusCode != fiber.StatusSeeOther {
		t.Fatalf("expected redirect status, got %d", resp.StatusCode)
	}
}

func TestProtectedRouteWithValidToken(t *testing.T) {
	middleware.SetJWTSecret("another-test-secret")
	tok, err := middleware.GenerateToken(55, "user")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	app := fiber.New()
	app.Get("/prot", middleware.AuthRequired, func(c *fiber.Ctx) error {
		uid := c.Locals("uid")
		if uid == nil {
			return c.SendStatus(500)
		}
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/prot", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: tok})

	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestValidateTokenRejectsInvalidToken(t *testing.T) {
	middleware.SetJWTSecret("secret-for-invalid-token")
	_, err := middleware.ValidateToken("this-is-not-a-valid-jwt")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestAuthStatusWithoutCookieDoesNotCrash(t *testing.T) {
	app := fiber.New()
	app.Get("/", middleware.AuthStatus, func(c *fiber.Ctx) error {
		if logged, ok := c.Locals("IsLogged").(bool); !ok || logged {
			return c.Status(500).SendString("expected not logged")
		}
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCheckSubscriptionAllowsActiveUser(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("UserRole", "user")
		c.Locals("uid", uint(1))
		c.Locals("Suscrito", true)
		c.Locals("FinPrueba", time.Now().AddDate(0, 0, 1))
		return c.Next()
	}, middleware.CheckSubscription, func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil), 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCheckSubscriptionRedirectsWhenExpired(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("UserRole", "user")
		c.Locals("uid", uint(1))
		c.Locals("Suscrito", false)
		c.Locals("FinPrueba", time.Now().AddDate(0, 0, -1))
		return c.Next()
	}, middleware.CheckSubscription, func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil), 5000)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if resp.StatusCode != fiber.StatusFound {
		t.Fatalf("expected redirect to /planes, got %d", resp.StatusCode)
	}
}
