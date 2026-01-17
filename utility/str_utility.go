package utility

import (
	"runtime"
	"strings"
)

func WinCharacter(s string) string {
	if runtime.GOOS == "windows" {
		// Caractères interdits sur Windows : < > : " / \ | ? *
		replacer := strings.NewReplacer(
			"<", "-",
			">", "-",
			":", "-",
			"\"", "'",
			"/", "-",
			"\\", "-",
			"|", "-",
			"?", "",
			"*", "-",
		)
		s = replacer.Replace(s)
		// Nettoyer les espaces multiples
		s = strings.ReplaceAll(s, "  ", " ")
		// Supprimer les espaces en début et fin
		s = strings.TrimSpace(s)
		// Supprimer les points et espaces en fin (interdits sur Windows)
		s = strings.TrimRight(s, ". ")
		return s
	}
	return s
}
