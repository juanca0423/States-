// Package utils
package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/template/handlebars/v3"
	"golang.org/x/exp/constraints"
)

var registerHelpersOnce sync.Once

func RegistrarHelpers(engine *handlebars.Engine) {
	registerHelpersOnce.Do(func() {
		engine.AddFunc("eq", func(a, b any) bool {
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		})

		engine.AddFunc("percent", func(current, total any) float64 {
			var curr, tot float64
			fmt.Sscanf(fmt.Sprintf("%v", current), "%f", &curr)
			fmt.Sscanf(fmt.Sprintf("%v", total), "%f", &tot)
			if tot == 0 {
				return 0
			}
			p := (curr / tot) * 100
			if p > 100 {
				return 100
			}
			return p
		})

		engine.AddFunc("lt", func(a, b int) bool {
			return a < b
		})

		engine.AddFunc("slice", func(s string, start, end int) string {
			if start < 0 {
				start = 0
			}
			if end > len(s) {
				end = len(s)
			}
			if start >= len(s) || start >= end {
				return ""
			}
			return s[start:end]
		})

		engine.AddFunc("fechaCorta", func(t time.Time) string {
			return t.Format("02/01/2006")
		})
	})
}

// max returns the greater of two ordered values.
func max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

