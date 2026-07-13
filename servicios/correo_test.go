package servicios

import (
	"strings"
	"testing"
)

func TestBuildVerificationLink(t *testing.T) {
	token := "abc123"
	got := buildVerificationLink(token)
	want := "https://eeffs.com/verificar?token=abc123"
	if got != want {
		t.Fatalf("buildVerificationLink() = %q; want %q", got, want)
	}
}

func TestBuildVerificationHtmlIncludesNameAndLink(t *testing.T) {
	html := buildVerificationHtml("Ana", "abc123")
	if !strings.Contains(html, "¡Hola Ana!") {
		t.Fatal("expected HTML to contain user name")
	}
	if !strings.Contains(html, "https://eeffs.com/verificar?token=abc123") {
		t.Fatal("expected HTML to contain verification link")
	}
}
