package config_test

import (
	"testing"

	"ef/config"
)

func TestFContFormatsPositiveValue(t *testing.T) {
    got := config.FCont(12345.678)
    want := "12,345.68"
    if got != want {
        t.Fatalf("FCont(12345.678) = %q; want %q", got, want)
    }
}

func TestFContFormatsNegativeValueWithParentheses(t *testing.T) {
    got := config.FCont(-1234.5)
    want := "(1,234.50)"
    if got != want {
        t.Fatalf("FCont(-1234.5) = %q; want %q", got, want)
    }
}

func TestFContZeroReturnsDash(t *testing.T) {
    got := config.FCont(0)
    want := "-"
    if got != want {
        t.Fatalf("FCont(0) = %q; want %q", got, want)
    }
}
