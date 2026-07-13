package utils_test

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ef/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v3"
)

func renderTemplate(t *testing.T, templateContent string, data any) string {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.hbs")
	if err := os.WriteFile(filePath, []byte(templateContent), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	engine := handlebars.New(dir, ".hbs")
	utils.RegistrarHelpers(engine)

	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("test", data)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	return string(body)
}

func TestEqHelper(t *testing.T) {
	got := renderTemplate(t, `{{#if (eq 1 1)}}OK{{else}}FAIL{{/if}}`, nil)
	if got != "OK" {
		t.Fatalf("expected OK, got %q", got)
	}
}

func TestPercentHelper(t *testing.T) {
	got := renderTemplate(t, `{{percent 2 4}}`, nil)
	if got != "50" {
		t.Fatalf("expected 50, got %q", got)
	}
}

func TestSliceHelper(t *testing.T) {
	got := renderTemplate(t, `{{slice "abcdef" 1 4}}`, nil)
	if got != "bcd" {
		t.Fatalf("expected bcd, got %q", got)
	}
}
