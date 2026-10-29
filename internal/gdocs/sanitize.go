package gdocs

import (
	"regexp"
	"strings"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"
)

var (
	// regexConsecutiveNewlines colapsa 3 ou mais quebras de linha consecutivas em no máximo duas (\n\n)
	regexConsecutiveNewlines = regexp.MustCompile(`\n{3,}`)
	// regexTrailingSpaces remove espaços em branco residuais no fim de cada linha
	regexTrailingSpaces = regexp.MustCompile(`[ \t]+\n`)
	// regexMultipleSpaces colapsa múltiplos espaços horizontais repetidos em um único espaço
	regexMultipleSpaces = regexp.MustCompile(`[ \t]{2,}`)
)

// SanitizeText aplica normalização canônica Unicode NFC e higieniza quebras de linha e espaços extras
func SanitizeText(raw string) string {
	if raw == "" {
		return ""
	}

	// 1. Normalização canônica Unicode NFC (converte acentos decompostos para pré-compostos)
	normalized := norm.NFC.String(raw)

	// 2. Padronização de quebras de linha (CRLF -> LF)
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	// 3. Remoção de espaços antes de quebra de linha
	normalized = regexTrailingSpaces.ReplaceAllString(normalized, "\n")

	// 4. Colapso de múltiplos espaços consecutivos (sem afetar novas linhas)
	normalized = regexMultipleSpaces.ReplaceAllString(normalized, " ")

	// 5. Colapso de quebras supérfluas (\n\n\n... -> \n\n)
	normalized = regexConsecutiveNewlines.ReplaceAllString(normalized, "\n\n")

	return normalized
}

// UTF16Length calcula com exatidão matemática o comprimento da string em code units UTF-16
func UTF16Length(s string) int64 {
	return int64(len(utf16.Encode([]rune(s))))
}

