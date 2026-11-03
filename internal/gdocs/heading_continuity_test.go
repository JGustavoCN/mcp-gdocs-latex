package gdocs

import (
	"regexp"
	"testing"
)

func TestAuditHeadingNumberingContinuity(t *testing.T) {
	// Caso com salto de 2.2 para 2.5 (esperado detectar ausência de 2.3 e 2.4)
	recordsWithGap := []HeadingRecord{
		{Text: "1 INTRODUÇÃO", Number: "1", Parts: []int{1}},
		{Text: "1.1 Contextualização", Number: "1.1", Parts: []int{1, 1}},
		{Text: "1.2 Problema", Number: "1.2", Parts: []int{1, 2}},
		{Text: "2 FUNDAMENTAÇÃO TEÓRICA", Number: "2", Parts: []int{2}},
		{Text: "2.1 Inteligência Artificial", Number: "2.1", Parts: []int{2, 1}},
		{Text: "2.2 Large Language Models", Number: "2.2", Parts: []int{2, 2}},
		{Text: "2.5 Arquitetura do Sistema", Number: "2.5", Parts: []int{2, 5}},
	}

	anomalies := AuditHeadingNumberingContinuity(recordsWithGap)
	if len(anomalies) != 1 {
		t.Fatalf("Esperava exatamente 1 anomalia de sequência, obteve %d", len(anomalies))
	}

	a := anomalies[0]
	if a.Actual != "2.5" || a.ParentPrefix != "2" {
		t.Errorf("Anomalia incorreta: %+v", a)
	}
	if a.Expected != "2.3, 2.4" {
		t.Errorf("Esperava ausentes '2.3, 2.4', obteve '%s'", a.Expected)
	}

	// Caso sequencial contínuo sem saltos
	continuousRecords := []HeadingRecord{
		{Text: "1 INTRODUÇÃO", Number: "1", Parts: []int{1}},
		{Text: "1.1 Contextualização", Number: "1.1", Parts: []int{1, 1}},
		{Text: "1.2 Problema", Number: "1.2", Parts: []int{1, 2}},
		{Text: "2 FUNDAMENTAÇÃO TEÓRICA", Number: "2", Parts: []int{2}},
		{Text: "2.1 Inteligência Artificial", Number: "2.1", Parts: []int{2, 1}},
		{Text: "2.2 Large Language Models", Number: "2.2", Parts: []int{2, 2}},
		{Text: "2.3 Retrieval-Augmented Generation", Number: "2.3", Parts: []int{2, 3}},
		{Text: "2.4 Model Context Protocol", Number: "2.4", Parts: []int{2, 4}},
		{Text: "2.5 Trabalhos Correlatos", Number: "2.5", Parts: []int{2, 5}},
	}

	cleanAnomalies := AuditHeadingNumberingContinuity(continuousRecords)
	if len(cleanAnomalies) != 0 {
		t.Errorf("Sequência contínua não deveria gerar anomalias, gerou: %+v", cleanAnomalies)
	}
}

func TestTableColumnSpecRegex(t *testing.T) {
	reColSpecPrefix := regexp.MustCompile(`^(?:(?:[lcrpXm]|p\{[^}]+\}|>\{[^}]+\}|@\{[^}]*\})\s*)+\}\s*`)

	cases := []struct {
		input       string
		expectedPre string
		expectedRem string
	}{
		{
			input:       "p{3.0cm} X X} Dimensão",
			expectedPre: "p{3.0cm} X X} ",
			expectedRem: "Dimensão",
		},
		{
			input:       "l c X} Categoria",
			expectedPre: "l c X} ",
			expectedRem: "Categoria",
		},
		{
			input:       "l c X} Prioridade",
			expectedPre: "l c X} ",
			expectedRem: "Prioridade",
		},
		{
			input:       "c c c X} Condição",
			expectedPre: "c c c X} ",
			expectedRem: "Condição",
		},
		{
			input:       ">{\\raggedright\\arraybackslash}p{3.2cm} X c c@{}} Módulo",
			expectedPre: ">{\\raggedright\\arraybackslash}p{3.2cm} X c c@{}} ",
			expectedRem: "Módulo",
		},
	}

	for _, c := range cases {
		loc := reColSpecPrefix.FindStringIndex(c.input)
		if loc == nil {
			t.Errorf("Falha ao casar especificação de coluna em %q", c.input)
			continue
		}
		matchedPre := c.input[:loc[1]]
		rem := c.input[loc[1]:]
		if matchedPre != c.expectedPre {
			t.Errorf("Prefixo incorreto em %q: obteve %q, esperava %q", c.input, matchedPre, c.expectedPre)
		}
		if rem != c.expectedRem {
			t.Errorf("Texto remanescente incorreto em %q: obteve %q, esperava %q", c.input, rem, c.expectedRem)
		}
	}
}

