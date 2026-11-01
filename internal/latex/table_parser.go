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
	regexTableEnv    = regexp.MustCompile(`(?s)\\begin\{(?:table|quadro)\}(.*?)\\end\{(?:table|quadro)\}`)
	regexCaption     = regexp.MustCompile(`\\caption\{([^}]+)\}`)
	regexLabel       = regexp.MustCompile(`\\label\{([^}]+)\}`)
	regexFonte       = regexp.MustCompile(`(?:Fonte:|\\fonte\{)([^}\n\\]+)`)
	regexTabular     = regexp.MustCompile(`(?s)\\begin\{(?:tabular|tabularx)\}(?:\{[^}]*\})*\{[^}]*\}(.*?)\\end\{(?:tabular|tabularx)\}`)
	regexTextCite    = regexp.MustCompile(`\\textcite\{([^}]+)\}`)
	regexCite        = regexp.MustCompile(`\\cite\{([^}]+)\}`)
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

	// 4. Extração do conteúdo dentro de tabular
	tabularContent := block
	if tabMatch := regexTabular.FindStringSubmatch(block); len(tabMatch) > 1 {
		tabularContent = tabMatch[1]
	}

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
		if trimmed == "" {
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

	// Remove \textbf{...}, \textit{...}, \texttt{...}, \path{...}, \url{...} preservando conteúdo interno
	boldRegex := regexp.MustCompile(`\\textbf\{([^}]+)\}`)
	s = boldRegex.ReplaceAllString(s, "$1")
	italicRegex := regexp.MustCompile(`\\textit\{([^}]+)\}`)
	s = italicRegex.ReplaceAllString(s, "$1")
	codeRegex := regexp.MustCompile(`\\(?:texttt|path|url)\{([^}]+)\}`)
	s = codeRegex.ReplaceAllString(s, "$1")
	refRegex := regexp.MustCompile(`\\(?:ref|pageref)\{([^}]+)\}`)
	s = refRegex.ReplaceAllString(s, "$1")

	// Colapso de múltiplos espaços
	s = strings.Join(strings.Fields(s), " ")

	return s
}
