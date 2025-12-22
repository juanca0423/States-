package config

import (
	"strconv"
	"strings"
)

func FCont(v float64) string {
	// 2 decimales, separador "." para decimales
	s := strconv.FormatFloat(v, 'f', 2, 64)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	// miles
	var out strings.Builder
	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out.WriteRune(',')
		}
		out.WriteRune(r)
	}
	out.WriteRune('.')
	out.WriteString(decPart)
	return out.String()
}
