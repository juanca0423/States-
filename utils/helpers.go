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

	engine.AddFunc("slice", func(s string, start, end int) string {
		if len(s) < end {
			return s
		}
		return s[start:end]
	})

	engine.AddFunc("fechaCorta", func(t time.Time) string {
		return t.Format("02/01/2006")
	})
}
