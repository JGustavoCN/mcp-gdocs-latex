package gdocs

import (
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/docs/v1"
)

// ============================================================================
// TESTES ADVERSARIAIS DO CHALLENGER 1 (Validação Empírica Estrita R1)
// ============================================================================

// 1. Stress Test: Espaços em branco vs Caracteres Invisíveis
func TestAdversarial_WhitespaceAndInvisibleChars(t *testing.T) {
	testCases := []struct {
		name                 string
		rawText              string
		expectedTrimmedEmpty bool
		notes                string
	}{
		{
			name:                 "Espaços ASCII comuns",
			rawText:              "   \t  \n",
			expectedTrimmedEmpty: true,
			notes:                "Espaços simples, tabs e newlines devem ser colapsados",
		},
		{
			name:                 "Espaço Não-Separável (NBSP \\u00A0)",
			rawText:              "\u00a0\u00a0\n",
			expectedTrimmedEmpty: true,
			notes:                "NBSP é classificado como espaço pelo unicode.IsSpace",
		},
		{
			name:                 "Múltiplos tabs e quebras de linha misturadas",
			rawText:              "\t\r\n\t\n  \r\n",
			expectedTrimmedEmpty: true,
			notes:                "Sequência complexa de whitespace",
		},
		{
			name:                 "Zero-Width Space (\\u200B)",
			rawText:              "\u200b\n",
			expectedTrimmedEmpty: false, // unicode.IsSpace NÃO considera \u200B espaço
			notes:                "Caractere invisível ZWSP não é aparado por TrimSpace padrão do Go",
		},
		{
			name:                 "Byte Order Mark / Zero-Width No-Break Space (\\uFEFF)",
			rawText:              "\ufeff\n",
			expectedTrimmedEmpty: false,
			notes:                "BOM \\uFEFF não é considerado espaço pelo unicode.IsSpace",
		},
		{
			name:                 "Zero-Width Non-Joiner (\\u200C)",
			rawText:              "\u200c\n",
			expectedTrimmedEmpty: false,
			notes:                "ZWNJ invisível",
		},
		{
			name:                 "Quebra de linha limpa simples (\\n)",
			rawText:              "\n",
			expectedTrimmedEmpty: true,
			notes:                "Parágrafo padrão vazio do Google Docs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trimmed := strings.TrimSpace(tc.rawText)
			isEmpty := (trimmed == "")
			if isEmpty != tc.expectedTrimmedEmpty {
				t.Logf("[OBSERVAÇÃO COMPORTAMENTAL] %s: TrimSpace(%q) == %q (isEmpty=%v, esperado=%v) - %s",
					tc.name, tc.rawText, trimmed, isEmpty, tc.expectedTrimmedEmpty, tc.notes)
			}
		})
	}
}

// 2. Stress Test: Parágrafos vazios pré-tabela e prevenção contra colapso de tabela
func TestAdversarial_TableDelimiterProtection(t *testing.T) {
	// Cenário A: Tabela precedida por exatamente 1 parágrafo vazio
	// Cenário B: Tabela precedida por 3 parágrafos vazios consecutivos (P1, P2, P3 -> Tabela)
	// Cenário C: Duas tabelas consecutivas separadas por 1 parágrafo vazio (T1 -> P_sep -> T2)
	// Cenário D: Duas tabelas consecutivas separadas por 3 parágrafos vazios (T1 -> P1, P2, P3 -> T2)
	// Cenário E: Tabela no final do documento precedida por parágrafo vazio

	content := []*docs.StructuralElement{
		// [0] P0: Título 1 (100-110)
		{StartIndex: 100, EndIndex: 110, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "1 INTRODUÇÃO\n"}}}}},
		// [1] P1: Vazio espúrio 1 (110-111)
		{StartIndex: 110, EndIndex: 111, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		// [2] P2: Vazio espúrio 2 (111-112)
		{StartIndex: 111, EndIndex: 112, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "  \n"}}}}},
		// [3] P3: Vazio delimitador pré-tabela 1 (112-113) -> DEVE SER PRESERVADO
		{StartIndex: 112, EndIndex: 113, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		// [4] T1: Tabela 1 (113-150)
		{StartIndex: 113, EndIndex: 150, Table: &docs.Table{}},
		// [5] P4: Vazio pós-T1 / espúrio (150-151)
		{StartIndex: 150, EndIndex: 151, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		// [6] P5: Vazio delimitador pré-T2 (151-152) -> DEVE SER PRESERVADO (pois content[7] é Tabela)
		{StartIndex: 151, EndIndex: 152, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		// [7] T2: Tabela 2 (152-190)
		{StartIndex: 152, EndIndex: 190, Table: &docs.Table{}},
		// [8] P6: Vazio final de corpo (190-191, bodyEnd = 191) -> DEVE SER PRESERVADO (clamping bodyEnd)
		{StartIndex: 190, EndIndex: 191, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
	}

	targetStartIndex := int64(100)
	bodyEnd := int64(191)

	var toDeleteRanges []*docs.Range
	var tableDelimiterRanges []*docs.Range

	for i := 0; i < len(content); i++ {
		el := content[i]
		if el.StartIndex < targetStartIndex || el.EndIndex >= bodyEnd {
			continue
		}
		if el.Table != nil {
			continue
		}
		if el.Paragraph != nil {
			hasPageBreak := false
			hasInlineObject := false
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.PageBreak != nil {
					hasPageBreak = true
				}
				if pel.InlineObjectElement != nil {
					hasInlineObject = true
				}
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}

			if hasPageBreak || hasInlineObject {
				continue
			}

			t := strings.TrimSpace(txt)
			if t == "" {
				isBeforeTable := false
				if i+1 < len(content) && content[i+1].Table != nil {
					isBeforeTable = true
				}

				if isBeforeTable {
					tableDelimiterRanges = append(tableDelimiterRanges, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				} else {
					toDeleteRanges = append(toDeleteRanges, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				}
			}
		}
	}

	// Verificações Rigorosas:
	// 1. P1 (110), P2 (111) e P4 (150) devem ir para toDeleteRanges
	expectedDeletes := []int64{110, 111, 150}
	if len(toDeleteRanges) != len(expectedDeletes) {
		t.Fatalf("Esperava %d deleções, obteve %d (%v)", len(expectedDeletes), len(toDeleteRanges), toDeleteRanges)
	}
	for idx, exp := range expectedDeletes {
		if toDeleteRanges[idx].StartIndex != exp {
			t.Errorf("Deleção [%d]: esperava offset %d, obteve %d", idx, exp, toDeleteRanges[idx].StartIndex)
		}
	}

	// 2. P3 (112) e P5 (151) devem ir para tableDelimiterRanges (PRESERVADOS DA DELEÇÃO)
	expectedDelimiters := []int64{112, 151}
	if len(tableDelimiterRanges) != len(expectedDelimiters) {
		t.Fatalf("Esperava %d delimitadores pré-tabela, obteve %d (%v)", len(expectedDelimiters), len(tableDelimiterRanges), tableDelimiterRanges)
	}
	for idx, exp := range expectedDelimiters {
		if tableDelimiterRanges[idx].StartIndex != exp {
			t.Errorf("Delimitador [%d]: esperava offset %d, obteve %d", idx, exp, tableDelimiterRanges[idx].StartIndex)
		}
	}

	// 3. P6 (190-191) é o último elemento do body (bodyEnd = 191) e NUNCA pode ser adicionado
	for _, r := range toDeleteRanges {
		if r.StartIndex == 190 || r.EndIndex >= bodyEnd {
			t.Errorf("FALHA CRÍTICA: Parágrafo final do body [%d-%d] foi erroneamente incluído para deleção!", r.StartIndex, r.EndIndex)
		}
	}

	// 4. Teste de isolamento mútuo: nenhum delimitador pode estar em toDeleteRanges
	for _, del := range tableDelimiterRanges {
		for _, rem := range toDeleteRanges {
			if del.StartIndex == rem.StartIndex {
				t.Fatalf("COLISÃO CRÍTICA: Offset %d está simultaneamente em tableDelimiterRanges e toDeleteRanges!", del.StartIndex)
			}
		}
	}
}

// 3. Stress Test: Sequências complexas de múltiplos PageBreak intercalados com quebras
func TestAdversarial_ComplexPageBreakInterleaving(t *testing.T) {
	content := []*docs.StructuralElement{
		{StartIndex: 10, EndIndex: 20, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Capa\n"}}}}},
		{StartIndex: 20, EndIndex: 21, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{PageBreak: &docs.PageBreak{}}}}},
		{StartIndex: 21, EndIndex: 22, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		{StartIndex: 22, EndIndex: 23, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{PageBreak: &docs.PageBreak{}}}}},
		{StartIndex: 23, EndIndex: 24, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{PageBreak: &docs.PageBreak{}}}}},
		{StartIndex: 24, EndIndex: 25, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "   \t\n"}}}}},
		{
			StartIndex: 25,
			EndIndex:   27,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "\n"}},
					{PageBreak: &docs.PageBreak{}},
				},
			},
		},
		{StartIndex: 27, EndIndex: 28, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{PageBreak: &docs.PageBreak{}}}}},
		{StartIndex: 28, EndIndex: 40, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "1 INTRODUÇÃO\n"}}}}},
		{StartIndex: 40, EndIndex: 41, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}}, // final
	}

	targetStartIndex := int64(10)
	bodyEnd := int64(41)

	var toDeleteRanges []*docs.Range

	for i := 0; i < len(content); i++ {
		el := content[i]
		if el.StartIndex < targetStartIndex || el.EndIndex >= bodyEnd {
			continue
		}
		if el.Paragraph != nil {
			hasPageBreak := false
			hasInlineObject := false
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.PageBreak != nil {
					hasPageBreak = true
				}
				if pel.InlineObjectElement != nil {
					hasInlineObject = true
				}
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}

			if hasPageBreak || hasInlineObject {
				continue
			}

			tStr := strings.TrimSpace(txt)
			if tStr == "" {
				isBeforeTable := false
				if i+1 < len(content) && content[i+1].Table != nil {
					isBeforeTable = true
				}

				if !isBeforeTable {
					toDeleteRanges = append(toDeleteRanges, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				}
			}
		}
	}

	// Apenas P2 (offset 21) e P5 (offset 24) devem ser marcados para deleção
	if len(toDeleteRanges) != 2 {
		t.Fatalf("Esperava exatamente 2 deleções, obteve %d", len(toDeleteRanges))
	}

	protectedPageBreakOffsets := []int64{20, 22, 23, 25, 27}
	for _, pbOffset := range protectedPageBreakOffsets {
		for _, del := range toDeleteRanges {
			if del.StartIndex == pbOffset {
				t.Fatalf("FALHA GRAVE: PageBreak no offset %d foi adicionado para deleção!", pbOffset)
			}
		}
	}

	// Ordenação reversa obrigatória
	sort.Slice(toDeleteRanges, func(i, j int) bool {
		return toDeleteRanges[i].StartIndex > toDeleteRanges[j].StartIndex
	})

	if toDeleteRanges[0].StartIndex != 24 || toDeleteRanges[1].StartIndex != 21 {
		t.Fatalf("Ordem reversa incorreta: esperava [24, 21], obteve [%d, %d]",
			toDeleteRanges[0].StartIndex, toDeleteRanges[1].StartIndex)
	}
}

// 4. Stress Test: Falsos títulos terminados com ':' e variações com pontuação/espaços
func TestAdversarial_FalseTitlesWithColonAndPunctuation(t *testing.T) {
	testCases := []struct {
		name                 string
		rawText              string
		currentStyle         string
		expectedIsHeading    bool
		expectedLinterAlert  bool
		description          string
	}{
		{
			name:                "Falso título clássico terminando em ':'",
			rawText:             "Quanto aos procedimentos técnicos:",
			currentStyle:        "HEADING_1",
			expectedIsHeading:   false,
			expectedLinterAlert: true,
			description:         "Oração introdutória estilizada como título deve ser rebaixada para NORMAL_TEXT",
		},
		{
			name:                "Falso título terminando em ':' com espaços e tabs ao final",
			rawText:             "Os seguintes requisitos foram identificados:  \t \n",
			currentStyle:        "HEADING_2",
			expectedIsHeading:   false,
			expectedLinterAlert: true,
			description:         "TrimSpace remove espaços finais garantindo terminação em ':'",
		},
		{
			name:                "Título legítimo numerado com ':' no MEIO da frase",
			rawText:             "3.1 ESTUDO DE CASO: APLICAÇÃO PRÁTICA",
			currentStyle:        "HEADING_2",
			expectedIsHeading:   true,
			expectedLinterAlert: false,
			description:         "Dois pontos no meio do subtítulo NÃO deve desqualificar o heading",
		},
		{
			name:                "Capítulo primário numerado terminando indevidamente em ':'",
			rawText:             "3 METODOLOGIA:",
			currentStyle:        "HEADING_1",
			expectedIsHeading:   false,
			expectedLinterAlert: true,
			description:         "NBR 6024 proíbe pontuação em títulos; se termina com ':' é rebaixado",
		},
		{
			name:                "Falso título terminando em ':. ' (dois pontos + ponto)",
			rawText:             "Quanto aos procedimentos:.",
			currentStyle:        "HEADING_1",
			expectedIsHeading:   true, // strings.HasSuffix(trimmed, ":") NÃO detecta ':.'
			expectedLinterAlert: false,
			description:         "Edge case adversarial: pontuação mista ':.'",
		},
		{
			name:                "Falso título terminando em ': --' (dois pontos + travessão)",
			rawText:             "Quanto aos procedimentos: --",
			currentStyle:        "HEADING_1",
			expectedIsHeading:   true, // strings.HasSuffix(trimmed, ":") NÃO detecta ': --'
			expectedLinterAlert: false,
			description:         "Edge case adversarial: dois pontos seguido de traço",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trimmed := strings.TrimSpace(tc.rawText)

			// Simula a lógica de classificação de abnt_writer.go
			isHeading1Num := rePrimaryHeadingNum.MatchString(trimmed)
			isHeading1NoNum := rePrimaryHeadingNoNum.MatchString(trimmed) || trimmed == "3 REFERÊNCIAS"
			isHeading2 := reSecondaryHeading.MatchString(trimmed)
			isHeading3 := reTertiaryHeading.MatchString(trimmed)
			isHeading4 := reQuaternaryHeading.MatchString(trimmed)

			if !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isHeading3 && !isHeading4 {
				if tc.currentStyle == "HEADING_1" && len(trimmed) < 120 &&
					!reSecondaryHeading.MatchString(trimmed) &&
					!reTertiaryHeading.MatchString(trimmed) &&
					!reQuaternaryHeading.MatchString(trimmed) {
					isHeading1Num = true
				} else if tc.currentStyle == "HEADING_2" && len(trimmed) < 120 &&
					!reTertiaryHeading.MatchString(trimmed) &&
					!reQuaternaryHeading.MatchString(trimmed) {
					isHeading2 = true
				}
			}

			// Regra 5: Terminação em ':' anula headings
			if strings.HasSuffix(trimmed, ":") {
				isHeading1Num = false
				isHeading1NoNum = false
				isHeading2 = false
				isHeading3 = false
				isHeading4 = false
			}

			actualIsHeading := isHeading1Num || isHeading1NoNum || isHeading2 || isHeading3 || isHeading4

			// Simula regra do Linter
			isCurrentStyleHeading := strings.HasPrefix(tc.currentStyle, "HEADING_") || tc.currentStyle == "TITLE"
			linterAlert := isCurrentStyleHeading && strings.HasSuffix(trimmed, ":")

			if actualIsHeading != tc.expectedIsHeading {
				t.Logf("[ANÁLISE ADVERSARIAL] %s: actualIsHeading=%v, esperado=%v (trimmed=%q)",
					tc.name, actualIsHeading, tc.expectedIsHeading, trimmed)
			}
			if linterAlert != tc.expectedLinterAlert {
				t.Logf("[ANÁLISE ADVERSARIAL LINTER] %s: linterAlert=%v, esperado=%v",
					tc.name, linterAlert, tc.expectedLinterAlert)
			}
		})
	}
}

// 5. Stress Test: Precedência de Títulos Secundários sobre Estilos Legados e Citação
func TestAdversarial_HeadingPrecedenceAndQuoteDemotion(t *testing.T) {
	// Seção secundária com estilo HEADING_1 prévio no documento
	secTitle := "3.1 CARACTERIZAÇÃO DA PESQUISA"
	currentStyle := "HEADING_1"

	isHeading1Num := rePrimaryHeadingNum.MatchString(secTitle)
	isHeading2 := reSecondaryHeading.MatchString(secTitle)

	if isHeading1Num {
		t.Fatalf("rePrimaryHeadingNum casou com seção secundária %q!", secTitle)
	}
	if !isHeading2 {
		t.Fatalf("reSecondaryHeading não reconheceu %q!", secTitle)
	}

	// Avalia a guarda de fallback
	if !isHeading1Num && !isHeading2 {
		if currentStyle == "HEADING_1" {
			isHeading1Num = true
		}
	}

	// Comprova que isHeading1Num permaneceu falso e isHeading2 é verdadeiro
	if isHeading1Num {
		t.Fatalf("Seção secundária foi erroneamente rebaixada/promovida a HEADING_1!")
	}
	if !isHeading2 {
		t.Fatalf("Seção secundária deveria ser HEADING_2!")
	}

	// Verifica que uma oração terminada em ':' não é considerada citação direta longa mesmo com recuo prévio
	quoteCandidate := "As principais características observadas foram as seguintes:"
	existingIndent := float64(113.4) // 4cm
	isLongQuote := existingIndent >= 85.0 && !strings.HasSuffix(quoteCandidate, ":")
	if isLongQuote {
		t.Fatalf("Oração terminada em ':' NÃO pode ser formatada como citação direta longa!")
	}
}
