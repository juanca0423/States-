// Package utils
package utils

import (
	"fmt"
	"time"

	"github.com/gofiber/template/handlebars/v3"
)

func RegistrarHelpers(engine *handlebars.Engine) {
	engine.AddFunc("eq", func(a, b any) bool {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	})

	engine.AddFunc("percent", func(current, total any) float64 {
		curr, _ := fmt.Sscanf(fmt.Sprintf("%v", current), "%f")
		tot, _ := fmt.Sscanf(fmt.Sprintf("%v", total), "%f")
		if tot == 0 {
			return 0
		}
		p := (float64(curr) / float64(tot)) * 100
		if p > 100 {
			return 100
		}
		return p
	})

	engine.AddFunc("lt", func(a, b int) bool {
		return a < b
	})

	engine.AddFunc("slice", func(s string, start, end int) string {
		if len(s) < end {
			return s
		}
		return s[start:end]
	})

	engine.AddFunc("eq", func(a, b any) bool {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	})

	engine.AddFunc("fechaCorta", func(t time.Time) string {
		return t.Format("02/01/2006")
	})
}
