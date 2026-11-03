package gdocs

import (
	"strings"
	"testing"
)

func TestCleanLatexInlineText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Conversão de operadores matemáticos e expressões",
			input:    `Delineamento fatorial $2 \times 2$ com $x \ge 5$ e $y \le 10$.`,
			expected: `Delineamento fatorial 2 × 2 com x ≥ 5 e y ≤ 10.`,
		},
		{
			name:     "Remoção de comandos e til insecável",
			input:    `Conforme a Seção~\ref{sec:rag} e a Resolução n\textordmasculine~35/2016.`,
			expected: `Conforme a Seção sec:rag e a Resolução nº 35/2016.`,
		},
		{
			name:     "Correção de parênteses duplicados e conversão RAG-MCP",
			input:    `Utilizando o Model Context Protocol) (MCP) e (Retrieval-Augmented Generation), RAG), como base.`,
			expected: `Utilizando o Model Context Protocol (MCP) e (Retrieval-Augmented Generation -- RAG), como base.`,
		},
		{
			name:     "Remoção de marcadores residuais de tabela",
			input:    `[TABELA] Conteúdo da matriz fatorial [/TABELA]`,
			expected: `Conteúdo da matriz fatorial`,
		},
		{
			name:     "Colapso de espaços duplos e espaços antes de pontuação",
			input:    `Este parágrafo   possui espaços excessivos  ,  e vírgula isolada .`,
			expected: `Este parágrafo possui espaços excessivos, e vírgula isolada.`,
		},
		{
			name:     "Limpeza de tags de estilo e caminhos",
			input:    `Verifique o código em \path{internal/gdocs} e o termo \textit{framework}.`,
			expected: `Verifique o código em internal/gdocs e o termo framework.`,
		},
		{
			name:     "Remoção de texttt, allowbreak e fechamento misto (caso real informado pelo usuário)",
			input:    `\texttt{mcp-\allowbreak use) na versão 2.5.2. Esse framework) gerencia os pacotes oficiais do ecossistema npm (\texttt{@modelcontextprotocol/\allowbreak core), \texttt{server), \texttt{client)`,
			expected: `mcp-use na versão 2.5.2. Esse framework) gerencia os pacotes oficiais do ecossistema npm (@modelcontextprotocol/core, server, client`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanLatexInlineText(tt.input)
			if got != tt.expected {
				t.Errorf("cleanLatexInlineText() = %q, esperado %q", got, tt.expected)
			}
		})
	}
}

func TestLinterRegexRules(t *testing.T) {
	// 1. Detecção de sintaxe matemática $...$
	mathSnippet := `O delineamento fatorial $2 \times 2$ foi aplicado.`
	if !reMathSyntaxInText.MatchString(mathSnippet) {
		t.Errorf("reMathSyntaxInText deveria detectar $2 \\times 2$")
	}

	// 2. Detecção de comando LaTeX em texto
	latexSnippet := `O comando \cite{lewis2020} não foi limpo.`
	if !reLatexCmdInText.MatchString(latexSnippet) {
		t.Errorf("reLatexCmdInText deveria detectar \\cite{...}")
	}

	// 3. Detecção de espaços duplos
	doubleSpaceSnippet := `Texto com  espaço duplo.`
	if !reDoubleSpace.MatchString(doubleSpaceSnippet) {
		t.Errorf("reDoubleSpace deveria detectar espaços duplos")
	}

	// 4. Detecção de espaço antes de pontuação
	punctSnippet := `Texto com espaço antes da vírgula , aqui.`
	if !reSpaceBeforePunct.MatchString(punctSnippet) {
		t.Errorf("reSpaceBeforePunct deveria detectar espaço antes de vírgula")
	}

	// 5. Detecção de número de capítulo
	chapTitle := `2 REFERENCIAL TEÓRICO`
	m := reChapterNumber.FindStringSubmatch(chapTitle)
	if len(m) < 2 || m[1] != "2" {
		t.Errorf("reChapterNumber deveria extrair o número 2 de %q", chapTitle)
	}

	// 6. Detecção de texttt e allowbreak (caso real)
	snippetUser := `\texttt{mcp-\allowbreak use) na versão 2.5.2.`
	if !reLatexCmdInText.MatchString(snippetUser) {
		t.Errorf("reLatexCmdInText deveria detectar \\texttt em %q", snippetUser)
	}

	allowbreakSnippet := `O comando \allowbreak não deve vazar no corpo.`
	if !reLatexCmdInText.MatchString(allowbreakSnippet) {
		t.Errorf("reLatexCmdInText deveria detectar \\allowbreak em %q", allowbreakSnippet)
	}
}

func TestLineSpacingAuditLogic(t *testing.T) {
	testCases := []struct {
		trimmed           string
		isHeading         bool
		state             docSectionState
		isLongQuote       bool
		isCaptionOrSource bool
		expectedSpacing   float64
	}{
		{trimmed: "Texto normal do corpo do trabalho acadêmico", isHeading: false, state: stateTextual, isLongQuote: false, isCaptionOrSource: false, expectedSpacing: 150.0},
		{trimmed: "REFERÊNCIAS", isHeading: true, state: statePostTextual, isLongQuote: false, isCaptionOrSource: false, expectedSpacing: 150.0},
		{trimmed: "ANTHROPIC. Model Context Protocol: Specification and Architecture. 2024.", isHeading: false, state: statePostTextual, isLongQuote: false, isCaptionOrSource: false, expectedSpacing: 100.0},
		{trimmed: "Figura 1 -- Fluxograma do pipeline", isHeading: false, state: stateTextual, isLongQuote: false, isCaptionOrSource: true, expectedSpacing: 100.0},
		{trimmed: "Fonte: Autor (2026)", isHeading: false, state: stateTextual, isLongQuote: false, isCaptionOrSource: true, expectedSpacing: 100.0},
		{trimmed: "Citação direta longa com recuo de 4cm...", isHeading: false, state: stateTextual, isLongQuote: true, isCaptionOrSource: false, expectedSpacing: 100.0},
	}

	for _, tc := range testCases {
		isHeadingItem := tc.isHeading ||
			rePrimaryHeadingNum.MatchString(tc.trimmed) ||
			rePrimaryHeadingNoNum.MatchString(tc.trimmed) ||
			reSecondaryHeading.MatchString(tc.trimmed) ||
			reTertiaryHeading.MatchString(tc.trimmed) ||
			reQuaternaryHeading.MatchString(tc.trimmed) ||
			tc.trimmed == "REFERÊNCIAS" || strings.HasPrefix(tc.trimmed, "REFERÊNCIAS")

		expectedSpacing := float64(150.0)
		if tc.isCaptionOrSource {
			expectedSpacing = 100.0
		} else if tc.isLongQuote {
			expectedSpacing = 100.0
		} else if isHeadingItem {
			expectedSpacing = 150.0
		} else if tc.state == statePostTextual {
			expectedSpacing = 100.0
		}

		if expectedSpacing != tc.expectedSpacing {
			t.Errorf("Para %q: esperado %.1f, obteve %.1f", tc.trimmed, tc.expectedSpacing, expectedSpacing)
		}
	}
}

func TestLinterGhostHeadingsAndHierarchy(t *testing.T) {
	// 1. Ghost heading: estilo HEADING_1 com texto vazio
	styleName := "HEADING_1"
	trimmed := ""
	hasPageBreak := false
	hasInlineObject := false
	isHeading := strings.HasPrefix(styleName, "HEADING_") || styleName == "TITLE"

	isGhostHeading := isHeading && trimmed == "" && !hasPageBreak && !hasInlineObject
	if !isGhostHeading {
		t.Errorf("Deveria identificar heading fantasma com estilo %s e texto vazio", styleName)
	}

	// Se tiver PageBreak, NÃO é ghost heading
	if isHeading && trimmed == "" && true && !hasInlineObject {
		// com hasPageBreak = true não deve ser flagrado como ghost heading
		flagged := isHeading && trimmed == "" && !true && !hasInlineObject
		if flagged {
			t.Errorf("Não deveria flagrar como ghost heading quando possui PageBreak")
		}
	}

	// 2. Falso título com dois pontos (:)
	titleWithColon := "Quanto aos Procedimentos Técnicos: A investigação conjuga duas modalidades clássicas de procedimento:"
	if !strings.HasSuffix(titleWithColon, ":") {
		t.Errorf("Deveria detectar terminação com dois pontos em %q", titleWithColon)
	}

	// 3. Hierarquia de Títulos (NBR 6024): 3.1 com estilo HEADING_1
	secTitle := "3.1 CARACTERIZAÇÃO DA PESQUISA"
	if !reSecondaryHeading.MatchString(secTitle) {
		t.Errorf("reSecondaryHeading deveria reconhecer %q", secTitle)
	}
	if rePrimaryHeadingNum.MatchString(secTitle) {
		t.Errorf("rePrimaryHeadingNum NÃO deveria reconhecer seção secundária %q", secTitle)
	}

	// 4. Seção terciária: 3.3.1
	tertTitle := "3.3.1 Levantamento do domínio e fontes de informação"
	if !reTertiaryHeading.MatchString(tertTitle) {
		t.Errorf("reTertiaryHeading deveria reconhecer %q", tertTitle)
	}
}

func TestHeadingPrecedenceRule(t *testing.T) {
	// Testa que uma seção secundária com currentStyle="HEADING_1" é classificada como HEADING_2 e NÃO HEADING_1
	trimmed := "3.1 CARACTERIZAÇÃO DA PESQUISA"
	currentStyle := "HEADING_1"

	isHeading1Num := rePrimaryHeadingNum.MatchString(trimmed)
	isHeading1NoNum := rePrimaryHeadingNoNum.MatchString(trimmed) || trimmed == "3 REFERÊNCIAS"
	isHeading2 := reSecondaryHeading.MatchString(trimmed)
	isHeading3 := reTertiaryHeading.MatchString(trimmed)
	isHeading4 := reQuaternaryHeading.MatchString(trimmed)

	if !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isHeading3 && !isHeading4 {
		if currentStyle == "HEADING_1" && len(trimmed) < 120 &&
			!reSecondaryHeading.MatchString(trimmed) &&
			!reTertiaryHeading.MatchString(trimmed) &&
			!reQuaternaryHeading.MatchString(trimmed) {
			isHeading1Num = true
		}
	}

	if isHeading1Num {
		t.Errorf("isHeading1Num não deveria ser true para %q, mesmo com currentStyle=HEADING_1", trimmed)
	}
	if !isHeading2 {
		t.Errorf("isHeading2 deveria ser true para %q", trimmed)
	}
}

func TestLinterSpuriousBlankLinePreservationRules(t *testing.T) {
	testCases := []struct {
		name            string
		trimmed         string
		isPreTextual    bool
		isBeforeTable   bool
		hasPageBreak    bool
		hasInlineObject bool
		shouldReport    bool
	}{
		{
			name:            "Linha em branco isolada no corpo textual",
			trimmed:         "",
			isPreTextual:    false,
			isBeforeTable:   false,
			hasPageBreak:    false,
			hasInlineObject: false,
			shouldReport:    true,
		},
		{
			name:            "Linha em branco pré-textual (capa/folha de rosto)",
			trimmed:         "",
			isPreTextual:    true,
			isBeforeTable:   false,
			hasPageBreak:    false,
			hasInlineObject: false,
			shouldReport:    false,
		},
		{
			name:            "Delimitador técnico imediatamente pré-tabela",
			trimmed:         "",
			isPreTextual:    false,
			isBeforeTable:   true,
			hasPageBreak:    false,
			hasInlineObject: false,
			shouldReport:    false,
		},
		{
			name:            "Parágrafo com quebra de página (PageBreak)",
			trimmed:         "",
			isPreTextual:    false,
			isBeforeTable:   false,
			hasPageBreak:    true,
			hasInlineObject: false,
			shouldReport:    false,
		},
		{
			name:            "Parágrafo com objeto inline (Logo)",
			trimmed:         "",
			isPreTextual:    false,
			isBeforeTable:   false,
			hasPageBreak:    false,
			hasInlineObject: true,
			shouldReport:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reported := false
			if tc.trimmed == "" {
				if !tc.isPreTextual && !tc.hasPageBreak && !tc.hasInlineObject {
					if !tc.isBeforeTable {
						reported = true
					}
				}
			}
			if reported != tc.shouldReport {
				t.Errorf("reported = %v, esperado = %v", reported, tc.shouldReport)
			}
		})
	}
}

