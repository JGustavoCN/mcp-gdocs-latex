package latex

import (
	"regexp"
	"strings"
)

// LaTeXTableData encapsula os dados e metadados de uma tabela extraída de código LaTeX
type LaTeXTableData struct {
	ID        string     `json:"id"`
	Label     string     `json:"label"`
	Caption   string     `json:"caption"`
	Source    string     `json:"source"`
	RowsCount int64      `json:"rows_count"`
	ColsCount int64      `json:"cols_count"`
	Headers   []string   `json:"headers"`
	Data      [][]string `json:"data"`
}

var (
	regexTableEnv = regexp.MustCompile(`(?s)\\begin\{(?:table|quadro)\}(.*?)\\end\{(?:table|quadro)\}`)
	regexCaption  = regexp.MustCompile(`\\caption\{([^}]+)\}`)
	regexLabel    = regexp.MustCompile(`\\label\{([^}]+)\}`)
	regexFonte    = regexp.MustCompile(`(?:Fonte:|\\fonte\{)([^}\n\\]+)`)
	regexTabular  = regexp.MustCompile(`(?s)\\begin\{(?:tabular|tabularx)\}(?:\{[^}]*\})*\{[^}]*\}(.*?)\\end\{(?:tabular|tabularx)\}`)
	regexTextCite = regexp.MustCompile(`\\textcite\{([^}]+)\}`)
	regexCite     = regexp.MustCompile(`\\cite\{([^}]+)\}`)
)

// ParseAllTablesInLatex extrai e normaliza todas as tabelas e quadros de um arquivo LaTeX
func ParseAllTablesInLatex(latexContent string) []LaTeXTableData {
	matches := regexTableEnv.FindAllStringSubmatch(latexContent, -1)
	var tables []LaTeXTableData

	for _, m := range matches {
		block := m[1]
		tbl := ParseSingleLaTeXTable(block)
		if tbl.RowsCount > 0 && tbl.ColsCount > 0 {
			tables = append(tables, tbl)
		}
	}

	return tables
}

// ParseSingleLaTeXTable analisa um bloco de tabela LaTeX e preenche a estrutura LaTeXTableData
func ParseSingleLaTeXTable(block string) LaTeXTableData {
	var tbl LaTeXTableData

	// 1. Extração de Caption
	if capMatch := regexCaption.FindStringSubmatch(block); len(capMatch) > 1 {
		tbl.Caption = CleanLaTeXText(capMatch[1])
	}

	// 2. Extração de Label
	if lblMatch := regexLabel.FindStringSubmatch(block); len(lblMatch) > 1 {
		tbl.Label = strings.TrimSpace(lblMatch[1])
		tbl.ID = tbl.Label
	}

	// 3. Extração de Fonte
	if fontMatch := regexFonte.FindStringSubmatch(block); len(fontMatch) > 1 {
		tbl.Source = CleanLaTeXText(fontMatch[1])
	} else if strings.Contains(block, "Autor (2026)") || strings.Contains(block, "autor (2026)") {
		tbl.Source = "Autor (2026)"
	}

	// 4. Extração do conteúdo dentro de tabular/tabularx com descarte robusto de colunas (balanceamento de chaves)
	tabularContent := extractTabularBody(block)

	// Remove regras de linha e espaçamentos do booktabs antes de particionar em linhas por \\
	cleanedTabular := tabularContent
	for _, cmd := range []string{"\\toprule", "\\midrule", "\\bottomrule", "\\hline", "\\addlinespace"} {
		cleanedTabular = strings.ReplaceAll(cleanedTabular, cmd, "")
	}

	// 5. Linhas do corpo da tabela
	rawLines := strings.Split(cleanedTabular, "\\\\")
	var grid [][]string

	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || isColumnSpecLine(trimmed) {
			continue
		}

		// Divide colunas por &
		rawCells := strings.Split(trimmed, "&")
		var row []string
		for _, cell := range rawCells {
			cleaned := CleanLaTeXText(cell)
			row = append(row, cleaned)
		}

		if len(row) > 0 {
			grid = append(grid, row)
		}
	}

	if len(grid) == 0 {
		return tbl
	}

	tbl.RowsCount = int64(len(grid))
	tbl.ColsCount = int64(len(grid[0]))

	tbl.Headers = grid[0]
	if len(grid) > 1 {
		tbl.Data = grid[1:]
	} else {
		tbl.Data = [][]string{}
	}

	return tbl
}

// CleanLaTeXText normaliza comandos específicos do LaTeX para texto legível
func CleanLaTeXText(raw string) string {
	s := strings.TrimSpace(raw)

	// Normaliza \textcite{key} -> key
	s = regexTextCite.ReplaceAllString(s, "$1")
	s = regexCite.ReplaceAllString(s, "($1)")

	// Desescape de %
	s = strings.ReplaceAll(s, `\%`, "%")
	s = strings.ReplaceAll(s, `\_`, "_")
	s = strings.ReplaceAll(s, `\&`, "&")

	// Expressões matemáticas inline e símbolos comuns
	s = strings.ReplaceAll(s, `\times`, "×")
	s = strings.ReplaceAll(s, `\ge`, "≥")
	s = strings.ReplaceAll(s, `\le`, "≤")
	s = strings.ReplaceAll(s, `\neq`, "≠")
	s = strings.ReplaceAll(s, `\textordmasculine`, "º")
	s = strings.ReplaceAll(s, `\textordfeminine`, "ª")
	s = strings.ReplaceAll(s, `\allowbreak`, "")
	s = strings.ReplaceAll(s, `~`, " ")
	s = strings.ReplaceAll(s, `$`, "")

	// Remove \texorpdfstring{TeX}{PDF} -> preserva argumento limpo (ex: 2 × 2)
	texOrPdfRegex := regexp.MustCompile(`\\texorpdfstring\{([^}]*)\}\{([^}]*)\}`)
	s = texOrPdfRegex.ReplaceAllString(s, "$1")

	// Remove \textbf{...}, \textit{...}, \texttt{...}, \path{...}, \url{...} preservando conteúdo interno
	boldRegex := regexp.MustCompile(`(?i)\\textbf\{([^}]+)\}`)
	s = boldRegex.ReplaceAllString(s, "$1")
	italicRegex := regexp.MustCompile(`(?i)\\textit\{([^}]+)\}`)
	s = italicRegex.ReplaceAllString(s, "$1")
	codeRegex := regexp.MustCompile(`(?i)\\(?:texttt|path|url)\{([^}]+)\}`)
	s = codeRegex.ReplaceAllString(s, "$1")
	refRegex := regexp.MustCompile(`(?i)\\(?:ref|pageref)\{([^}]+)\}`)
	s = refRegex.ReplaceAllString(s, "$1")

	// Remove resíduos órfãos de comandos de formatação que sobraram sem fechar chave
	s = regexp.MustCompile(`(?i)\\(?:textit|textbf|texttt|path|url|texorpdfstring)\b`).ReplaceAllString(s, "")

	// Remove comandos residuais de alinhamento de colunas LaTeX
	colCmdsRegex := regexp.MustCompile(`(?:>\{[^}]*\}|@\{[^}]*\}|\\arraybackslash|\\raggedright|\\raggedleft|\\centering)`)
	s = colCmdsRegex.ReplaceAllString(s, "")

	// Colapso de múltiplos espaços
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// extractTabularBody extrai o conteúdo interno de tabular ou tabularx descartando
// integralmente a especificação de colunas, mesmo quando contiver chaves aninhadas (ex: p{3.6cm}, @{}).
func extractTabularBody(block string) string {
	startIdx := strings.Index(block, `\begin{tabular`)
	if startIdx == -1 {
		return block
	}

	isTabularX := strings.HasPrefix(block[startIdx:], `\begin{tabularx}`)
	var afterBegin string
	if isTabularX {
		afterBegin = block[startIdx+len(`\begin{tabularx}`):]
	} else {
		endBegin := strings.Index(block[startIdx:], "}")
		if endBegin == -1 {
			return block
		}
		afterBegin = block[startIdx+endBegin+1:]
	}

	// tabularx pula 2 blocos de chaves: {largura} e {colunas}
	// tabular pula 1 bloco de chaves: {colunas}
	numArgsToSkip := 1
	if isTabularX || strings.Contains(block[startIdx:startIdx+20], "tabular*") {
		numArgsToSkip = 2
	}

	cursor := 0
	for arg := 0; arg < numArgsToSkip; arg++ {
		openBrace := strings.Index(afterBegin[cursor:], "{")
		if openBrace == -1 {
			break
		}
		cursor += openBrace
		depth := 0
		closed := false
		for i := cursor; i < len(afterBegin); i++ {
			if afterBegin[i] == '{' {
				depth++
			} else if afterBegin[i] == '}' {
				depth--
				if depth == 0 {
					cursor = i + 1
					closed = true
					break
				}
			}
		}
		if !closed {
			break
		}
	}

	body := afterBegin[cursor:]
	endIdx := strings.Index(body, `\end{tabular`)
	if endIdx != -1 {
		body = body[:endIdx]
	}

	return body
}

// isColumnSpecLine verifica se uma linha crua é resquício da especificação de colunas do tabular
func isColumnSpecLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if strings.Contains(trimmed, `\arraybackslash`) || strings.Contains(trimmed, `\raggedright`) ||
		strings.Contains(trimmed, `\raggedleft`) || strings.Contains(trimmed, `\centering`) ||
		strings.Contains(trimmed, `@{}`) || strings.Contains(trimmed, `@{`) ||
		strings.HasPrefix(trimmed, "@{") || strings.HasSuffix(trimmed, "@{}}") ||
		strings.HasPrefix(trimmed, "l >{") || strings.HasPrefix(trimmed, "c >{") ||
		strings.HasPrefix(trimmed, "r >{") || strings.HasPrefix(trimmed, ">{\\raggedright") ||
		strings.HasPrefix(trimmed, "l c c") || strings.HasPrefix(trimmed, "l l l") ||
		strings.HasPrefix(trimmed, "c c c") || strings.HasPrefix(trimmed, "r r r") ||
		regexp.MustCompile(`^(?:l|c|r|p\{|X)\s+`).MatchString(trimmed) ||
		regexp.MustCompile(`>[{\\]`).MatchString(trimmed) {
		return true
	}
	return false
}
