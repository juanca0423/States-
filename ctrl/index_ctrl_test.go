package ctrl_test

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ef/ctrl"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v3"
)

func TestGetIndexRouteRenders(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "index.hbs"), []byte("INDEX_OK"), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	engine := handlebars.New(tmp, ".hbs")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", ctrl.GetIndex)

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil), 5000)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if string(body) != "INDEX_OK" {
		t.Fatalf("expected INDEX_OK, got %q", string(body))
	}
}

func TestGetLoguinRouteRenders(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "loguin.hbs"), []byte("LOGIN_OK"), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	engine := handlebars.New(tmp, ".hbs")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/loguin", ctrl.GetLoguin)

	resp, err := app.Test(httptest.NewRequest("GET", "/loguin", nil), 5000)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if string(body) != "LOGIN_OK" {
		t.Fatalf("expected LOGIN_OK, got %q", string(body))
	}
}
