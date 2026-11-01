package gdocs

import (
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

