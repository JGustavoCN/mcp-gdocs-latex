package gdocs

import (
	"testing"
)

func TestUTF16Length(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected int64
	}{
		{"ASCII", "Hello World", 11},
		{"Acentos em Português", "Atenção à validação", 19},
		{"Emoji (Surrogate Pair)", "Documento 🚀", 12}, // "Documento " (10) + emoji (2) = 12
		{"Vazio", "", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := UTF16Length(tc.input)
			if got != tc.expected {
				t.Errorf("UTF16Length(%q) = %d; esperado %d", tc.input, got, tc.expected)
			}
		})
	}
}

func TestSanitizeText(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Colapso de quebras consecutivas",
			input:    "Parágrafo 1\n\n\n\n\nParágrafo 2",
			expected: "Parágrafo 1\n\nParágrafo 2",
		},
		{
			name:     "Remoção de espaços no final de linha",
			input:    "Linha com espaços no fim    \nOutra linha",
			expected: "Linha com espaços no fim\nOutra linha",
		},
		{
			name:     "Colapso de múltiplos espaços no texto",
			input:    "Texto   com    muitos      espaços.",
			expected: "Texto com muitos espaços.",
		},
		{
			name:     "Normalização CRLF para LF",
			input:    "Linha 1\r\nLinha 2\r\n",
			expected: "Linha 1\nLinha 2\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeText(tc.input)
			if got != tc.expected {
				t.Errorf("SanitizeText(%q) = %q; esperado %q", tc.input, got, tc.expected)
			}
		})
	}
}

